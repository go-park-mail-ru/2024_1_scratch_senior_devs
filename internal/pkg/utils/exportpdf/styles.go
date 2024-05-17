package exportpdf

const stylesDark = `
* {
    box-sizing: border-box;
}

body, html {
    margin: 0;
    padding: 0;
}

body {
    background: #1b2735;
    color: #f4f6f9;
    font-family: sans-serif;
    font-size: 16px;
}

[data-selected='true']::marker {
    content: '☑ ';
}

[data-selected='false']::marker {
    content: '☐ ';
}

[data-selected='true'] {
    text-decoration: line-through;
}

[data-selected='false'] {
    text-decoration: none;
}

.note-editor-content {
    display: flex;
    flex-direction: column;
}

.note-title {
    width: 100%;
    padding: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
    font-weight: bold;
    font-size: 36px;
    border-bottom: 2px solid #2b3b55;
    margin: 0;
}

.note-body {
    padding: 20px;
}

.note-body > div {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
}

.attach-wrapper {
    outline: none;
    background: #18232f;
    padding: 8px 16px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    gap: 10px;
    text-decoration: none;
    position: relative;
    cursor: pointer;
    transition: 0.2s ease;
    margin-right: 2px;
    border: 1px solid #2b3b55;
}

.file-extension-label {
    background: #546189;
    color: #fff;
    border-radius: 8px;
    height: 25px;
    padding: 4px;
    display: flex;
    justify-content: center;
    align-items: center;
}

.file-name {
    color: #f4f6f9;
}

.close-attach-btn-container {
    background: #546189;
    padding: 8px;
    border-radius: 50%;
    transition: 0.3s;
    content: "x";
}

.close-attach-btn {
    filter: invert(80%);
    width: 20px;
    height: 20px;
}

.subnote-wrapper {
    background: #18232f;
    padding: 8px 16px;
    border-radius: 16px;
    transition: 0.2s ease;
    border: 1px solid #2b3b55;
    margin: 15px 2px;
    user-select: none;
    text-decoration: none;
}

.subnote-title {
    color: #f4f6f9;
}

.delete-subnote-btn-container {
    display: none;
}

.close-attach-btn-container {
    display: none;
}
`

const stylesLight = `
* {
    box-sizing: border-box;
}

body, html {
    margin: 0;
    padding: 0;
}

body {
    background: #fff;
    color: #0C090A;
    font-family: sans-serif;
    font-size: 16px;
}

[data-selected='true']::marker {
    content: '☑ ';
}

[data-selected='false']::marker {
    content: '☐ ';
}

[data-selected='true'] {
    text-decoration: line-through;
}

[data-selected='false'] {
    text-decoration: none;
}

.note-editor-content {
    display: flex;
    flex-direction: column;
}

.note-title {
    width: 100%;
    padding: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
    font-weight: bold;
    font-size: 36px;
    border-bottom: 2px solid #2b3b55;
    margin: 0;
}

.note-body {
    padding: 20px;
}

.note-body > div {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
}

.attach-wrapper {
    outline: none;
    background: #18232f;
    padding: 8px 16px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    gap: 10px;
    text-decoration: none;
    position: relative;
    cursor: pointer;
    transition: 0.2s ease;
    margin-right: 2px;
    border: 1px solid #2b3b55;
}

.file-extension-label {
    background: #546189;
    color: #fff;
    border-radius: 8px;
    height: 25px;
    padding: 4px;
    display: flex;
    justify-content: center;
    align-items: center;
}

.file-name {
    color: #f4f6f9;
}

.close-attach-btn-container {
    background: #546189;
    padding: 8px;
    border-radius: 50%;
    transition: 0.3s;
    content: "x";
}

.close-attach-btn {
    filter: invert(80%);
    width: 20px;
    height: 20px;
}

.subnote-wrapper {
    background: #18232f;
    padding: 8px 16px;
    border-radius: 16px;
    transition: 0.2s ease;
    border: 1px solid #2b3b55;
    margin: 15px 2px;
    user-select: none;
    text-decoration: none;
}

.subnote-title {
    color: #f4f6f9;
}

.delete-subnote-btn-container {
    display: none;
}

.close-attach-btn-container {
    display: none;
}
`
