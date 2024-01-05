$(function () {
  console.log('[gistrunner]: inject.js is loaded.');

  async function getGistInfo() {
    const gistId = window.location.pathname.split('/')[2];
    try {
      return await $.ajax({
        url: `https://api.github.com/gists/${gistId}`,
        async: false,
      });
    } catch (error) {
      console.error(`[gistrunner]: An error occurred while retrieving gist info ${gistId}`);
      console.error(error);
      return;
    }
  }
  function browserName() {
    let browserName;

    if (navigator.userAgent.match(/chrome|chromium|crios/i)) {
      browserName = 'Chrome';
    } else if (navigator.userAgent.match(/firefox|fxios/i)) {
      browserName = 'Firefox';
    } else if (navigator.userAgent.match(/safari/i)) {
      browserName = 'Safari';
    } else if (navigator.userAgent.match(/opr/ / i)) {
      browserName = 'Opera';
    } else if (navigator.userAgent.match(/edg/i)) {
      browserName = 'Edge';
    } else {
      browserName = '';
    }

    return browserName;
  }

  function createNewIssueUrl(rawGistUrl, language, responseBody) {
    const url = new URL('https://github.com/prettyirrelevant/gistrunner/issues/new');
    url.searchParams.set('title', `${browserName()} Extension Error`);
    url.searchParams.set(
      'body',
      `## Language\n${language}\n\n## Response\n${responseBody}\n\n## Gist URL\n${rawGistUrl}\n\n## Browser Info\n${navigator.userAgent}`
    );

    return url;
  }

  function fileNameToCssId(fileName) {
    const transformedString = fileName.replace(/[^a-zA-Z0-9-_]+/g, '-');
    return `file-${transformedString.replace(/^-+|-+$/g, '')}`;
  }

  async function getSupportedLanguages() {
    const { baseServerURL, supportedLanguages } = await chrome.storage.local.get([
      'baseServerURL',
      'supportedLanguages',
    ]);

    // TODO(prettyirrelevant): set expiry to invalidate local storage cache.
    if (supportedLanguages) return JSON.parse(supportedLanguages);

    try {
      result = await $.ajax({
        url: `${baseServerURL}/api/languages`,
        async: false,
      });
      if (!result.data) return;

      chrome.storage.local
        .set({ supportedLanguages: JSON.stringify(result.data) })
        .then(() => console.log('[gistrunner]: supported languages stored!'));

      return result.data;
    } catch (error) {
      console.error('[gistrunner]: An error occurred while fetching supported languages');
      console.error(error);
      return;
    }
  }

  async function insertRunButtonForFiles() {
    const gistInfo = await getGistInfo();
    if (!gistInfo) {
      return;
    }

    const supportedLanguages = await getSupportedLanguages();
    if (!supportedLanguages) {
      return;
    }

    const extensionToLanguageMap = {};
    Object.entries(supportedLanguages).forEach(([language, extension]) => {
      extensionToLanguageMap[extension] = language;
    });

    Object.values(gistInfo.files).forEach((el) => {
      const fileExt = `.${el.filename.split('.').slice(-1)}`;

      // TODO(prettyirrelevant): find a way to handle truncated files.
      const language = extensionToLanguageMap[fileExt];
      if (!el.truncated && language) {
        const fileActionDiv = $(`div#${fileNameToCssId(el.filename)} > div.file-header > div.file-actions`);
        const runButton = $('<a>', {
          class: 'Button--primary Button--small Button',
          style: 'margin-bottom:0px',
          on: {
            click: function (event) {
              handleRunClick(event, el, language);
            },
          },
          append: $('<span>', {
            class: 'Button-content',
            append: $('<span>', {
              class: 'Button-label',
              text: 'Run',
            }),
          }),
        });

        runButton.appendTo(fileActionDiv);
      }
    });
  }

  async function handleRunClick(event, fileInfo, language) {
    const { baseServerURL } = await chrome.storage.local.get(['baseServerURL']);
    const outputDivClassName = `${fileNameToCssId(fileInfo.filename)}-output`;

    $(`.${outputDivClassName}`).remove();
    $(event.target).find(`span.Button-label`).text('Running');
    $(event.target).prop('disabled', true);

    $.ajax({
      type: 'POST',
      async: false,
      url: `${baseServerURL}/api/run`,
      contentType: 'application/json',
      data: JSON.stringify({ language: language, content: fileInfo.content }),
      success: function (data, textStatus, jqXHR) {
        $(`div#${fileNameToCssId(fileInfo.filename)}`).after(
          `
              <div class="file mt-1 mb-5 ${outputDivClassName}">
                  <div class="file-header">
                      <div class="file-info">Output: ${fileInfo.filename}</div>
                  </div>

                  <div itemprop="text" class="Box-body p-2 blob-wrapper data type-log gist-border-0" style="font-family: ui-monospace,SFMono-Regular,SF Mono,Menlo,Consolas,Liberation Mono,monospace;font-size:0.9em">
                      ${data.data.Result.split('\n')
                        .map((line) => `<div>${line}</div>`)
                        .join('')}
                  </div>
              </div>
          `
        );

        // scroll to that point on the page
        $('html').animate({ scrollTop: $(`div.${outputDivClassName}`).offset().top }, 800);
      },
      error: function (jqXHR) {
        $(`div#${fileNameToCssId(fileInfo.filename)}`).after(
          `
                <div class="file mt-1 mb-5 ${outputDivClassName}">
                    <div class="file-header">
                        <div class="file-info">Output: ${fileInfo.filename}</div>
                    </div>

                    <div itemprop="text" class="Box-body p-2 blob-wrapper data type-log gist-border-0" style="font-family: ui-monospace,SFMono-Regular,SF Mono,Menlo,Consolas,Liberation Mono,monospace;font-size:0.9em">
                        <code>
                            Could not run ${
                              fileInfo.filename
                            }. Check console for more details. More details below↓ <br><br>
                            <pre>${jqXHR.responseText}</pre>
                            <a target="_blank" href=${createNewIssueUrl(
                              fileInfo.raw_url,
                              language,
                              jqXHR.responseText
                            )}>Open an issue</a>
                        </code>
                    </div>
                </div>
            `
        );

        // scroll to that point on the page
        $('html').animate({ scrollTop: $(`div.${outputDivClassName}`).offset().top }, 800);
      },
      complete: function (jqXHR, textStatus) {
        $(event.target).find(`span.Button-label`).text('Run');
        $(event.target).prop('disabled', false);
      },
    });
  }

  // entrypoint
  insertRunButtonForFiles();
});
