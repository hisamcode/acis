
function createSpanHtmxIndicator() {
  var spanHtmxIndicator = document.createElement("span")
  spanHtmxIndicator.classList.add("loading", "loading-spinner", "htmx-indicator")
  return spanHtmxIndicator
}
// function createSpanHtmxIndicatorDots() {
//   var spanHtmxIndicator = document.createElement("span")
//   spanHtmxIndicator.classList.add("loading", "loading-dots", "htmx-indicator", "loading-xs")
//   return spanHtmxIndicator
// }

// function createSkeletonHtmxIndicator() {
//   var span = document.createElement("span")
//   span.classList.add("")

//   var flex = document.createElement("div")
//   flex.classList.add("flex", "flex-col", "gap-4", "w-full", "my-4")

//   var skeleton = document.createElement("div")
//   skeleton.classList.add("skeleton", "h-4", "w-full")

//   flex.appendChild(skeleton)
//   span.appendChild(flex)

//   return span
// }

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

  document.body.addEventListener("clearValidation", function (evt) {
    let validationErrors = document.getElementsByClassName("validation_error")
    for (let i = 0; i < validationErrors.length; i++) {
      let validationError = validationErrors[i]
      let span = validationError.getElementsByTagName("span")
      if (span.length > 0) {
        span[0].remove()
      }
    }

  })

  document.body.addEventListener("toastCreateSuccess", function (evt) {
    generateToast("Create Success")
  })
  document.body.addEventListener("toastUpdateSuccess", function (evt) {
    console.log(evt)
    generateToast("Update Success")
  })
  document.body.addEventListener("toastDeleteSuccess", function (evt) {
    generateToast("Delete Success")
  })
  // var elementHtmxLoading = content.getElementsByClassName("insert-htmx-loading-dots")
  // if (elementHtmxLoading.length > 0) {
  //   var spanhtmxindicator = createSpanHtmxIndicatorDots()
  //   for (let i = 0; i < elementHtmxLoading.length; i++) {
  //     elementHtmxLoading[i].appendChild(spanhtmxindicator)
  //     console.log(elementHtmxLoading[i])
  //   }
  // }
  // var elementSkeleton = content.getElementsByClassName("insert-htmx-skeleton")
  // if (elementSkeleton.length > 0) {
  //   var elementHtmxLoading = createSkeletonHtmxIndicator()
  //   for (let i = 0; i < elementSkeleton.length; i++) {
  //     elementSkeleton[i].appendChild(elementHtmxLoading)
  //   }
  // }
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

document.addEventListener('htmx:afterRequest', function (evt) {
  if (evt.detail.xhr.status == 404) {
    /* Notify the user of a 404 Not Found response */
    return alert("Error: Could Not Find Resource");
  }
  if (evt.detail.successful != true) {
    /* Notify of an unexpected error, & print error to console */
    alert("Unexpected Error");
    return console.error(evt);
  }
  if (evt.detail.target.id == 'list_emoji') {
    if (evt.detail.successful == true) {
      form_emoji_create.reset()
      emoji_modal.close()
      emoji_modal_delete.close()
    }
  }
});


