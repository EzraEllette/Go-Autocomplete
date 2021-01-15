function getCaret(textarea) {
  if (textarea.selectionStart) {
      return textarea.selectionStart;
  } else if (document.selection) {
      textarea.focus();

      var range = document.selection.createRange();
      if (range == null) {
          return 0;
      }

      var re = textarea.createTextRange(),
          rc = re.duplicate();
      re.moveToBookmark(r.getBookmark());
      rc.setEndPoint('EndToStart', re);

      return rc.text.length;
  }
  return 0;
}

function InsertText(text) {
  var textarea = document.getElementById('textarea');
  var currentPos = getCaret(textarea);
  //alert(currentPos);
  var strLeft = textarea.value.substring(0, currentPos);
  var strMiddle = text;
  var strRight = textarea.value.substring(currentPos, textarea.value.length);
  textarea.value = strLeft + strMiddle + strRight;
}

function getWord() {
  const textarea = document.getElementById('textarea');
  let caret = getCaret(document.getElementById('textarea')) - 1;
  let textValue = textarea.value;

  let word = ""

  while (caret >= 0 && /\S/.test(textValue[caret])) {
    word = textValue[caret] + word
    caret--;
  }

  return word;
}

document.addEventListener('DOMContentLoaded', () => {
  // This will wait for the astilectron namespace to be ready
  document.addEventListener('astilectron-ready', function() {
    // This will send a message to GO
    const sendMessage = () => {
      let word = getWord().toLocaleLowerCase();
      if (word.length === 0) return;
      astilectron.sendMessage(word, function(message) {
        InsertText(message + ' ')
      });
    };

    document.getElementById('suggest-completions').addEventListener('click', sendMessage);
  });
});