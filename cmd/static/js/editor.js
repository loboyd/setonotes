// initialize the editor
var easyMDE = new EasyMDE({ element: document.getElementById("input-body") });
easyMDE.toggleFullScreen();
easyMDE.togglePreview();
easyMDE.isEditor = false;

// focus the preview div so the user can scroll with the keyboard immediately
previewDiv = document.getElementsByClassName("editor-preview-full")[0]
previewDiv.tabIndex='0';
setTimeout(function(){previewDiv.focus();},10);

// add a title field to the editor
var editorArea = document.getElementsByClassName('CodeMirror-sizer')[0];
var editorTitle = document.createElement('input');
editorTitle.value = document.getElementById('input-title').value;
editorTitle.id = 'editor-title';
editorTitle.autocomplete = 'off';
editorTitle.type = 'text';
editorArea.insertBefore(editorTitle, editorArea.childNodes[0]);

// capture keystrokes for editor control
document.onkeydown = function(evt) {
    evt = evt || window.event;
    //console.log(evt);
    // if the user presses `i` and the editor isn't active, then launch the editor
    if (evt.key == 'i' || evt.key == 'a' && !evt.ctrlKey && !evt.shiftKey &&
	    !easyMDE.isEditor) {
        evt.preventDefault();
        easyMDE.togglePreview();
        easyMDE.isEditor = true;
    // if the user presses `ESC` and the editor is active, then collapse the editor
    } else if (evt.key == 'Escape' && !evt.ctrlKey && !evt.shiftKey &&
	    easyMDE.isEditor) {
        easyMDE.togglePreview();
        easyMDE.isEditor = false;
    // if the user presses CTRL+Enter, then save and exit
    } else if (evt.key == 'Enter' && evt.ctrlKey && !evt.shiftKey) {
        easyMDE.saveNote();
    // `j` scrolls down in preview mode
    } else if (evt.key == 'j' && !evt.ctrlKey && !evt.shiftKey) {
        evt.preventDefault();
        previewDiv.scrollBy(0, 100);
    // 'k' scrolls up in preview mode
    } else if (evt.key == 'k' && !evt.ctrlKey && !evt.shiftKey) {
        evt.preventDefault();
        previewDiv.scrollBy(0,-100);
    }
};
