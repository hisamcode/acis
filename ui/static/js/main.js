
function createSpanHtmxIndicator() {
  var spanHtmxIndicator = document.createElement("span")
  spanHtmxIndicator.classList.add("loading", "loading-spinner", "htmx-indicator")
  return spanHtmxIndicator
}

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

  // var elementSkeleton = content.getElementsByClassName("insert-htmx-skeleton")
  // if (elementSkeleton.length > 0) {
  //   var elementHtmxLoading = createSkeletonHtmxIndicator()
  //   for (let i = 0; i < elementSkeleton.length; i++) {
  //     elementSkeleton[i].appendChild(elementHtmxLoading)
  //   }
  // }
})
