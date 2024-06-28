# ROADMAP

### v0.1 - Prototype

Building the simplest thing to show to my friends. Uses as few dependencies as possible, global state variables, etc. Built on free tier Gemini Flash.

#### tasks

- [x] start actually recording my tasks instead of winging it 
- [x] setup free tier gemini project to replace paid tier
- [x] pass history as context to query
- [ ] handle error responses
  - [error codes](https://docs.gemini.com/rest-api/#error-codes)
  - [error payloads](https://docs.gemini.com/rest-api/#error-payload)
- [x] layout
- [x] style message boxes (optional)
- [x] contact/attribution/legal in footer
- [x] tip jar
- [x] simple hard-coded links for loading different cats
- [ ] publish v0.1 to google, run server

## future ideas

- feature: visual interactivity
  - modify the cat image using the Vertex API and Imagen, based on the question and answer
- feature: share link to thread
- debt: organize stateful code in javascript (using React/Angular/etc)
- debt: install lodash or equivalent; replace imperative ops + bad code style
- feature: put more effort into visual experience
- 