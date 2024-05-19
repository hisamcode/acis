
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

htmx.onLoad(function (content) {
  var elementHtmxLoading = content.getElementsByClassName("insert-htmx-loading")
  if (elementHtmxLoading.length > 0) {
    var spanhtmxindicator = createSpanHtmxIndicator()
    for (let i = 0; i < elementHtmxLoading.length; i++) {
      elementHtmxLoading[i].appendChild(spanhtmxindicator)
    }
  }

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
