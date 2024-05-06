

htmx.onLoad(function (content) {
  var spanHtmxIndicator = document.createElement("span")
  spanHtmxIndicator.classList.add("loading", "loading-spinner", "htmx-indicator")

  var btns = content.getElementsByClassName("insert-htmx-loading")
  for (let i = 0; i < btns.length; i++) {
    btns[i].appendChild(spanHtmxIndicator)
  }
})
