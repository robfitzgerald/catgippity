const form = document.getElementById('catQueryForm');
const textField = document.getElementById('catQueryInput');
const responseParagraph = document.getElementById('response');
const catImageDiv = document.getElementById('catImageDiv');
const responseDiv = document.getElementById('responseDiv');
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
  responseDiv.appendChild(spinner)
  isLoading = true;
}

function removeSpinner() {
  responseDiv.removeChild(responseDiv.lastChild)
  isLoading = false;
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

      const catWelcome = document.createElement("p")
      catWelcome.innerText = responseText
      responseDiv.appendChild(catWelcome)
      // responseParagraph.innerText = responseText
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
  const thisQuestion = document.createElement("p")
  thisQuestion.innerText = textField.value
  responseDiv.appendChild(thisQuestion)

  appendSpinner()

  const formData = {
    question: textField.value,
    history: ""
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
      const responseText = data.cat_talk
      const catAdvice = document.createElement("p")
      catAdvice.innerText = responseText
      responseDiv.appendChild(catAdvice)

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