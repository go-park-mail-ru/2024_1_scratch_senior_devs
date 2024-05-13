package exportpdf

//const example = `
//<div class="note-editor-content">
//	<div class="note-title" contenteditable="true" data-placeholder="">Список покупок</div>
//	<div class="note-body">
//		<div contenteditable="true">
//			<ul data-type="todo">
//				<li data-selected="false">огурцы</li>
//				<li data-selected="false">хлеб</li>
//				<li data-selected="false">йогурты</li>
//				<li data-selected="false">шапка</li>
//				<li data-selected="false">куртка</li>
//				<li data-selected="false">табуретка</li>
//				<li data-selected="false">
//					<b><font color="#e00000" data-cursordayakrut="4">танк</font></b>
//				</li>
//			</ul>
//			<img contenteditable="false" width="500" src="https://you-note.ru/images/6a34815e-1ce9-47c1-b5b1-6862c59072d0.webp" /><br />
//			<a href="https://you-note.ru/notes/8f6621e9-68ba-468a-907f-5fcf0c8ae9ae">
//				<div class="subnote-container">
//					<img src="https://you-note.ru/src/assets/note.svg" class="subnote-icon" />
//					<span class="subnote-title">Подзаметка</span>
//					<div class="delete-subnote-btn-container">
//						<img src="https://you-note.ru/src/assets/trash.svg" class="delete-subnote-btn" />
//					</div>
//				</div>
//			</a>
//			<br />
//			<div><br /></div>
//		</div>
//	</div>
//</div>
//`

//func TestProcessButton(t *testing.T) {
//	input := `
//		<button class="subnote-wrapper" contenteditable="false" data-noteid="8f6621e9-68ba-468a-907f-5fcf0c8ae9ae">
//			<div class="subnote-container">
//				<span class="subnote-title">Подзаметка</span>
//				<div class="delete-subnote-btn-container"></div>
//			</div>
//		</button>
//	`
//	document, _ := goquery.NewDocumentFromReader(strings.NewReader(input))
//	processButton(document)
//	html, _ := document.Html()
//	fmt.Println(html)
//}

//func TestProcessIframe(t *testing.T) {
//	input := `
//		<iframe contenteditable="false" src="https://www.youtube.com/embed/U9Is9s9Ewgg"></iframe>
//	`
//	document, _ := goquery.NewDocumentFromReader(strings.NewReader(input))
//	processIframe(document)
//	html, _ := document.Html()
//	fmt.Println(html)
//}

//func TestPrepareHTML(t *testing.T) {
//	testCases := []struct {
//		name     string
//		input    string
//		expected string
//	}{
//		{
//			name:     "TestPrepareHTML_Success",
//			input:    `<img src="blob:https://you-note.ru/d4bbe6ee-fe5a-47cf-845a-f3f9f6869016" data-imgid="aa498c1c-04b2-4047-8814-6caa0a131237" />`,
//			expected: `<html><head></head><body><img src="https://you-note.ru/attaches/aa498c1c-04b2-4047-8814-6caa0a131237.webp" data-imgid="aa498c1c-04b2-4047-8814-6caa0a131237"/></body></html>`,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			output, err := PrepareHTML(tc.input)
//			if err != nil {
//				t.Errorf(err.Error())
//			}
//
//			assert.Equal(t, tc.expected, output)
//		})
//	}
//}

//func TestPrepareHTML(t *testing.T) {
//	input := `
//		<div class="note-editor-content">
//			<div class="note-title" contenteditable="true" data-placeholder="">Список покупок</div>
//			<div class="note-body">
//				<div contenteditable="true">
//					<ul data-type="todo">
//						<li data-selected="false">esfgогурцы<br /></li>
//						<li data-selected="false">хлеб</li>
//						<li data-selected="false">йогурты</li>
//						<li data-selected="false">шапка</li>
//						<li data-selected="false">куртка</li>
//						<li data-selected="false">табуретка</li>
//						<li data-selected="false">
//							<b><font color="#e00000">танк</font></b>
//						</li>
//					</ul>
//					<h1>header 1</h1>
//					<h2>header 2</h2>
//					<h2>header 3</h2>
//					<button class="attach-wrapper" contenteditable="false" data-fileid="6d774c99-e60d-400f-b82b-12622645093b" data-filename="Требования к отчёту по ЛР(СПО_ИУ5).pdf">
//						<div class="attach-container">
//							<div class="file-extension-label">pdf</div>
//							<span class="file-name">Требования к отчё...</span>
//							<div class="close-attach-btn-container"><img src="./src/assets/close.svg" class="close-attach-btn" /></div>
//						</div>
//					</button>
//					<ul>
//						<li>awdgsjhfk</li>
//						<li>dsfgh</li>
//						<li>sdfg</li>
//					</ul>
//					<ol>
//						<li>wert</li>
//						<li>34ty</li>
//						<ul>
//							<li>ewrer</li>
//							<li>ewfg</li>
//							<li><br /></li>
//						</ul>
//					</ol>
//					<div><br /></div>
//					<img contenteditable="false" width="500" src="blob:https://you-note.ru/4c217f66-d2a4-46fd-bc98-cf176119c9f1" data-imgid="aa498c1c-04b2-4047-8814-6caa0a131237" />
//					<div><br /></div>
//					<div><i>dsf</i>g 3<b>455</b></div>
//					<br />
//					<div><br /></div>
//					<div>
//						<u>
//							dfs&nbsp;
//							<s>
//								weg s<font color="#e00000">gdhf&nbsp;<span style="background-color: rgb(92, 75, 26);"> ghj</span></font>
//							</s>
//						</u>
//					</div>
//					<br />
//					<br />
//					<button class="subnote-wrapper" contenteditable="false" data-noteid="8f6621e9-68ba-468a-907f-5fcf0c8ae9ae">
//						<div class="subnote-container">
//							<img src="./src/assets/note.svg" class="subnote-icon" /><span class="subnote-title">Подзаметка</span>
//							<div class="delete-subnote-btn-container"><img src="./src/assets/trash.svg" class="delete-subnote-btn" /></div>
//						</div>
//					</button>
//					<br />
//					<div><br /></div>
//					<iframe contenteditable="false" src="https://www.youtube.com/embed/mUBXUyRoQco"></iframe><br />
//					<br />
//					<br />
//					<div>hello</div>
//					<div class="block-chosen blockplaceholder" data-cursordayakrut="0"><br /></div>
//					<br />
//					<div><br /></div>
//				</div>
//			</div>
//		</div>
//	`
//	result, err := prepareHTML(input)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	fmt.Println(result)
//}
