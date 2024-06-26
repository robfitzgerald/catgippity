// HTML object identifiers
const form = document.getElementById('catQueryForm');
const textField = document.getElementById('catQueryInput');
const catImageDiv = document.getElementById('catImageDiv');
const conversationDiv = document.getElementById('responseDiv');

// speaker attribute values
const catSpeaker = "cat"
const userSpeaker = "user"

// URLs
const apiUrl = '/welcome';
const welcomeUrl = '/welcome';
const queryUrl = '/query';
const catIds = {
  "red": 0,
  "dick": 1
}
var isLoading = false;

function appendSpinner() {
  const spinner = document.createElement("md-circular-progress")
  spinner.setAttribute("indeterminate", "")
  conversationDiv.appendChild(spinner)
  isLoading = true;
}

function removeSpinner() {
  conversationDiv.removeChild(conversationDiv.lastChild)
  isLoading = false;
}

function addToHistory(speaker, content) {
  const entry = document.createElement("p")
  entry.innerText = content
  entry.setAttribute("speaker", speaker)
  conversationDiv.appendChild(entry)
}

function getHistoryAsText() {
  var history = ""
  for (const element of conversationDiv.children) {
    const text = element.innerText
    const speaker = element.getAttribute("speaker")
    history += speaker + ": " + text + "\n"
  }
  // conversationDiv.children.forEach(element => {
  //   const text = element.innerText
  //   const speaker = element.getAttribute("speaker")
  //   history += speaker + ": " + text + "\n"
  // });
  return history
}

function welcome(cat_name) {
  appendSpinner()
  console.log("welcome called with cat: " + cat_name)
  const catId = catIds[cat_name]
  // const formData = {
  //   cat_id: cat_id
  // };
  const welcomeUrlWithCatId = welcomeUrl + "/" + catId

  fetch(welcomeUrlWithCatId, {
    method: 'GET',
    // headers: { 'Content-Type': 'application/json' },
    // body: JSON.stringify(formData)
  })
    .then(response => response.json())
    .then(data => {
      console.log(data);
      const responseText = data.cat_talk
      const responseImgUrl = data.image_url
      while (catImageDiv.firstChild) {
        catImageDiv.removeChild(catImageDiv.firstChild);
      }

      removeSpinner()

      const catImg = document.createElement("img");
      catImg.src = responseImgUrl;
      catImageDiv.appendChild(catImg)

      addToHistory(catSpeaker, responseText)
      // const catWelcome = document.createElement("p")
      // catWelcome.innerText = responseText
      // responseDiv.appendChild(catWelcome)
    })
    .catch(error => {
      console.error('Error:', error);
      if (isLoading) {
        removeSpinner()
      }
    })
    .finally(() => {
      if (isLoading) {
        removeSpinner()
      }
    })

}

form.addEventListener('submit', event => {

  event.preventDefault();

  // add this question to the queue
  addToHistory(userSpeaker, textField.value)
  // const thisQuestion = document.createElement("p")
  // thisQuestion.innerText = textField.value
  // responseDiv.appendChild(thisQuestion)

  appendSpinner()

  const formData = {
    question: textField.value,
    history: getHistoryAsText()
  };

  fetch(queryUrl, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData)
  })
    .then(response => response.json())
    .then(data => {
      console.log(data);
      removeSpinner()
      // const responseText = data.cat_talk
      // const catAdvice = document.createElement("p")
      // catAdvice.innerText = responseText
      // responseDiv.appendChild(catAdvice)
      addToHistory(catSpeaker, data.cat_talk)

    })
    .catch(error => {
      console.error('Error:', error);
      if (isLoading) {
        removeSpinner()
      }
    })
    .finally(() => {
      if (isLoading) {
        removeSpinner()
      }
    })
});