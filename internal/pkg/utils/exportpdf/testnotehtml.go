package exportpdf

const TestNoteHTMLInput = `
<div class="note-editor-content">
    <div class="note-title" contenteditable="true" data-placeholder="">Список покупок</div>
    <div class="note-body">
        <div id="note-editor-inner" contenteditable="true">
            <ul data-type="todo">
                <li data-selected="false">огурцы</li>
                <li data-selected="false">хлеб</li>
                <li data-selected="true" data-cursordayakrut3d22178dc9f2-421c-9faf-382d7adcbec4="0">йогурты</li>
                <li data-selected="true">шапка</li>
                <li data-selected="false">куртка</li>
                <li data-selected="false">табуретка</li>
                <li data-selected="false">
                    <b><font color="#e00000">танк</font></b>
                </li>
            </ul>
            <h1>header 1</h1>
            <h2>header 2</h2>
            <h2>header 3</h2>
            <button class="attach-wrapper" contenteditable="false" data-fileid="6d774c99-e60d-400f-b82b-12622645093b" data-filename="Требования к отчёту по ЛР(СПО_ИУ5).pdf">
                <div class="attach-container">
                    <div class="file-extension-label">pdf</div>
                    <span class="file-name">Требования к отчё...</span>
                    <div class="close-attach-btn-container"><img src="./src/assets/close.svg" class="close-attach-btn" /></div>
                </div>
            </button>
            <ul>
                <li><font color="#eab308">опять</font></li>
                <li data-cursordayakrut01696d3606b4-4af6-9d46-1650beb814f1="3"><s>эти</s></li>
                <li><span style="background-color: rgb(26, 92, 32);">списки</span></li>
            </ul>
            <ol>
                <li>
                    и <u>нумерованные</u> <font color="#9333ea"><b>тоже</b></font>
                </li>
                <li>ага</li>
                <ul>
                    <li>и <font color="#ffa500">ещё</font> <span style="background-color: rgb(92, 26, 58);">такие</span></li>
                    <li><i>миллион</i> <span style="background-color: rgb(92, 58, 26);">списков</span> короче</li>
                    <li>
                        <font color="#e00000"><b>!!!</b></font>
                    </li>
                </ul>
            </ol>
            <div><br /></div>
            <img contenteditable="false" data-imgid="aa498c1c-04b2-4047-8814-6caa0a131237" class="img" src="blob:https://you-note.ru/987baf89-1f9c-49e4-9ba4-ecaaf372cf81" />
            <div><br /></div>
            <div><i>а тут просто текстик такой простенький</i></div>
            <button class="subnote-wrapper" contenteditable="false" data-noteid="8f6621e9-68ba-468a-907f-5fcf0c8ae9ae">
                <div class="subnote-container">
                    <img src="./src/assets/note.svg" class="subnote-icon" /><span class="subnote-title">Подзаметка</span>
                    <div class="delete-subnote-btn-container"><img src="./src/assets/trash.svg" class="delete-subnote-btn" /></div>
                </div>
            </button>
            <div><br /></div>
            <div>и ещё можно даже видосы с ютуба сюда закидывать</div>
            <div><br /></div>
            <iframe contenteditable="false" src="https://www.youtube.com/embed/dcm8D8NVbmU"></iframe><br />
            <div><br /></div>
        </div>
    </div>
</div>
`
