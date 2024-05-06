

htmx.onLoad(function (content) {
  function makeSpanHtmxIndicator() {
    var span = document.createElement("span")
    span.classList.add("loading", "loading-spinner", "htmx-indicator")
    return span
  }
  var btns = content.getElementsByClassName("insert-htmx-loading")
  var spanHtmxIndicator = makeSpanHtmxIndicator()
  console.log(btns)
  for (let i = 0; i < btns.length; i++) {
    btns[i].appendChild(spanHtmxIndicator)
  }
})
