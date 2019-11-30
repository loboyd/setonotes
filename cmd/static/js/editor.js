// initialize the editor
var easyMDE = new EasyMDE({ element: document.getElementById("input-body") });
var initialEditorValue = easyMDE.value();
easyMDE.toggleFullScreen();

// add a title field to the editor
var editorArea = document.getElementsByClassName('CodeMirror-sizer')[0].firstChild;
var editorTitle = document.createElement('input');
editorTitle.value = document.getElementById('input-title').value;
editorTitle.id = 'editor-title';
editorTitle.autocomplete = 'off';
editorTitle.type = 'text';
editorTitle.onclick = function(){ this.focus(); }; // override codemirror focus
editorArea.insertBefore(editorTitle, editorArea.childNodes[0]);

// launch into preview mode
easyMDE.togglePreview();

// focus the preview div so the user can scroll with the keyboard immediately
previewDiv = document.getElementsByClassName("editor-preview-full")[0];
previewDiv.tabIndex='0';
setTimeout(function(){previewDiv.focus();},10);

mathDelimiters = [
    {left: "$$",  right: "$$",  display: true},
    {left: "\\[", right: "\\]", display: true},
    {left: "$",   right: "$",   display: false},
    {left: "\\(", right: "\\)", display: false},
];
renderMathInElement(document.body, {delimiters: mathDelimiters});

function escapeMarkdown(unescapedMarkdown) {
    escapedMarkdown = unescapedMarkdown.replace('*','\*');

    return escapedMarkdown;
}

// capture keystrokes for editor control
document.onkeydown = function(evt) {
    evt = evt || window.event;
    //console.log(evt);
    // if the user presses `i` and the editor isn't active, then launch the editor
    if ((evt.key == 'i' || evt.key == 'a') && !evt.ctrlKey && !evt.shiftKey &&
	    easyMDE.isPreviewActive()) {
        evt.preventDefault();
        easyMDE.togglePreview();
    // if the user presses `ESC` and the editor is active, then collapse the editor
    } else if (evt.key == 'Escape' && !evt.ctrlKey && !evt.shiftKey &&
	   !easyMDE.isPreviewActive()) {
        easyMDE.togglePreview();
	// this math rendering should really go inside togglePreview(), but I
        // couldn't figure out how to do that easily, so for now it lives
        // here TODO
        renderMathInElement(document.body, {delimiters: mathDelimiters});
    // if the user presses CTRL+Enter, then save and exit
    } else if (evt.key == 'Enter' && evt.ctrlKey && !evt.shiftKey) {
        easyMDE.saveNote();
    // `j` scrolls down in preview mode
    } else if (evt.key == 'j' && !evt.ctrlKey && !evt.shiftKey &&
	    easyMDE.isPreviewActive()) {
        evt.preventDefault();
        previewDiv.scrollBy(0, 100);
    // 'k' scrolls up in preview mode
    } else if (evt.key == 'k' && !evt.ctrlKey && !evt.shiftKey &&
	    easyMDE.isPreviewActive()) {
        evt.preventDefault();
        previewDiv.scrollBy(0,-100);
    // 'Backspace' goes back to the directory in preview mode
    } else if (evt.key == 'Backspace' && !evt.ctrlKey &&
	    !evt.shiftKey && easyMDE.isPreviewActive()) {
        evt.preventDefault();
        window.location.href = '../../';
    }
};

// if the user tries to leave without saving, prompt them
window.onbeforeunload = function(evt) {
    if (easyMDE.value() != initialEditorValue) {
        return true;
    }
}
