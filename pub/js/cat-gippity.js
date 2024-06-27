// HTML object identifiers
const form = document.getElementById('catQueryForm');
const textField = document.getElementById('catQueryInput');
const catImageDiv = document.getElementById('catImageDiv');
const conversationList = document.getElementById('conversationList');

// speaker attribute values
const catSpeaker = "cat"
const userSpeaker = "user"

// URLs
const apiUrl = '/welcome';
const welcomeUrl = '/welcome';
const queryUrl = '/query';
const catIds = new Map([
  ["red", 0],
  ["dick", 1]
])
var isLoading = false;
var catName = "dick"

function getCatFromUrl() {
  const queryString = window.location.search;
  console.log(queryString); // Outputs the entire query string (e.g., ?param1=value1&param2=value2)

  // Accessing specific parameters
  const urlParams = new URLSearchParams(queryString);
  const catIdLookup = urlParams.get('cat_id');
  const catNameSelected = catIds.get(catIdLookup)
  if (catNameSelected != undefined) {
    catName = catNameSelected
  }
}

function activateSpinner() {
  const spinner = document.createElement("md-circular-progress")
  spinner.setAttribute("indeterminate", "")
  conversationList.appendChild(spinner)
  isLoading = true;
}

function deactivateSpinner() {
  conversationList.removeChild(conversationList.lastChild)
  isLoading = false;
}

function addToHistory(speaker, content) {
  //       <md-card class="user-message">
  //   <div class="message-content">
  //   <p>This is a message from the user.</p>
  // </div>
  // </md-card>
  const card = document.createElement("md-card")
  card.className = "user-message"
  card.setAttribute("speaker", speaker)

  const textContent = document.createElement("p")
  textContent.className = "chat-content"
  textContent.innerText = content
  card.appendChild(textContent)

  // const entry = document.createElement("ul")
  // entry.className = "mdc-list mdc-list--two-line"
  // entry.innerText = content
  // entry.setAttribute("speaker", speaker)
  conversationList.appendChild(card)
}

function getHistoryAsText() {
  var history = ""
  for (const element of conversationList.children) {
    const text = element.firstChild.innerText
    const speaker = element.getAttribute("speaker")
    history += speaker + ": " + text + "\n"
  }

  return history
}

function popQuestion() {
  const question = textField.value.trim()
  textField.value = ""
  return question
}

function welcome() {
  activateSpinner()
  console.log("welcome called with cat: " + catName)
  const catId = catIds.get(catName)
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

      deactivateSpinner()

      const catImg = document.createElement("img");
      catImg.className = "chat-image"
      catImg.src = responseImgUrl;
      catImageDiv.appendChild(catImg)

      addToHistory(catSpeaker, responseText)
    })
    .catch(error => {
      console.error('Error:', error);
      if (isLoading) {
        deactivateSpinner()
      }
    })
    .finally(() => {
      if (isLoading) {
        deactivateSpinner()
      }
    })

}

form.addEventListener('submit', event => {

  event.preventDefault();
  question = popQuestion()
  if (question != "") {
    const history = getHistoryAsText()
    const formData = {
      question: question,
      history: history,
    };

    // update the page
    addToHistory(userSpeaker, question)
    activateSpinner()

    console.log("sending request:")
    console.log(formData)

    fetch(queryUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formData)
    })
      .then(response => response.json())
      .then(data => {
        console.log(data);
        deactivateSpinner()

        addToHistory(catSpeaker, data.cat_talk)

      })
      .catch(error => {
        console.error('Error:', error);
        if (isLoading) {
          deactivateSpinner()
        }
      })
      .finally(() => {
        if (isLoading) {
          deactivateSpinner()
        }
      })
  }
});

// Handle Enter key press in text field (optional)
textField.addEventListener('keydown', (event) => {
  if (event.key === 'Enter') {
    form.requestSubmit(); // Trigger form submission manually
  }
});

