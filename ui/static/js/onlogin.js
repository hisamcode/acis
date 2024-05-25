
htmx.onLoad(function (content) {
  console.log("on load ok")
  if (typeof (list_transactions) != 'undefined') {
    setAttributeClientDate(list_transactions)
  }
});

function onSubmit(evt) {
  let event = new CustomEvent("clearValidation")
  document.dispatchEvent(event)
}

(function () {
  document.addEventListener("toastCreateSuccess", function (evt) {
    generateToast("Create Success")
  })
  document.addEventListener("toastUpdateSuccess", function (evt) {
    generateToast("Update Success")
  })
  document.addEventListener("toastDeleteSuccess", function (evt) {
    generateToast("Delete Success")
  })


  document.addEventListener("clearValidation", function (evt) {
    let validationErrors = document.getElementsByClassName("validation_error")
    for (let i = 0; i < validationErrors.length; i++) {
      let validationError = validationErrors[i]
      let span = validationError.getElementsByTagName("span")
      if (span.length > 0) {
        span[0].remove()
      }
    }

  })

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

})()
