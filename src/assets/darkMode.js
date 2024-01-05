const toggleSwitch = document.querySelector("#darkModeToggle")
const isDarkMode = localStorage.getItem("darkMode")

if (isDarkMode === "true") {
  enableDarkMode()
  if (toggleSwitch) {
    toggleSwitch.checked = true
  }
}

if (toggleSwitch) {
  toggleSwitch.addEventListener("change", function () {
    if (this.checked) {
      enableDarkMode()
      localStorage.setItem("darkMode", "true")
    } else {
      disableDarkMode()
      localStorage.setItem("darkMode", "false")
    }
  })
}

function enableDarkMode() {
  document.body.classList.add("dark")
  document.body.style.backgroundColor = "#121212"
}

function disableDarkMode() {
  document.body.classList.remove("dark")
  document.body.style.backgroundColor = "#ffffff"
}
