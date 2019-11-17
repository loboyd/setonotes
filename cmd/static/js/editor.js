// initialize the editor
var easyMDE = new EasyMDE({ element: document.getElementById("input-body") });
easyMDE.toggleFullScreen();
easyMDE.togglePreview();
easyMDE.isEditor = false;
document.getElementsByClassName("editor-preview-full")[0].tabIndex='0';
setInterval(
    function() {
        document.getElementsByClassName("editor-preview-full")[0].focus();
    }, 10);

// capture keystrokes for editor control
document.onkeydown = function(evt) {
    evt = evt || window.event;
    console.log(evt);
    // if the user presses `i` and the editor isn't active, then launch the editor
    if (evt.keyCode == 73 && !easyMDE.isEditor) {
        evt.preventDefault();
        easyMDE.togglePreview();
        easyMDE.isEditor = true;
        easyMDE.codemirror.focus();
        // TODO set cursor position to top of preview window (and scroll, too)
        //easyMDE.codemirror.setCursor({line: 1, ch: 1}); // something like this
    // if the user presses `ESC` and the editor is active, then collapse the editor
    } else if (evt.keyCode == 27 && easyMDE.isEditor) {
        easyMDE.togglePreview();
        easyMDE.isEditor = false;
    // if the user presses CTRL+Enter, then save and exit
    } else if (evt.keyCode == 13 && evt.ctrlKey) {
        document.getElementById("input-form").submit();
    }
};
