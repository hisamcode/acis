
function createSpanHtmxIndicator() {
  var spanHtmxIndicator = document.createElement("span")
  spanHtmxIndicator.classList.add("loading", "loading-spinner", "htmx-indicator")
  return spanHtmxIndicator
}

function setAttributeClientDate(el) {
  let dateNow = new Date().toISOString()
  el.setAttribute("hx-headers", `{"client-date": "${dateNow}"}`)

  htmx.process(el)
}

// feedback for creating resource
function generateToast(string) {
  let div = document.createElement("div")
  div.classList.add("toast", "toast-top", "toast-center")
  div.id = "success_create"
  let divAlert = document.createElement("div")
  divAlert.classList.add("alert", "alert-info", "flex")
  let span = document.createElement("span")
  span.innerHTML = string
  let btn = document.createElement("button")
  btn.classList.add("btn", "btn-sm")
  btn.innerHTML = "close"
  btn.onclick = function () {
    div.remove()
  }

  divAlert.appendChild(span)
  divAlert.appendChild(btn)
  div.appendChild(divAlert)

  document.body.appendChild(div)
}


htmx.onLoad(function (content) {
  var elementHtmxLoading = content.getElementsByClassName("insert-htmx-loading")
  if (elementHtmxLoading.length > 0) {
    var spanhtmxindicator = createSpanHtmxIndicator()
    for (let i = 0; i < elementHtmxLoading.length; i++) {
      elementHtmxLoading[i].appendChild(spanhtmxindicator)
    }
  }

})

// timeServerToClient is time from server to current client time
// timeStr "09 May 24 08:22 UTC" or time from server with UTC
function timeServerToClient(timeStr, el) {
  var getLanguage = () => navigator.userLanguage || (navigator.languages && navigator.languages.length && navigator.languages[0]) || navigator.language || navigator.browserLanguage || navigator.systemLanguage || 'en';
  var timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
  var date = new Date(timeStr).toLocaleString(getLanguage(), {
    timeZone: timezone,
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  })
  el.innerHTML = date
}



