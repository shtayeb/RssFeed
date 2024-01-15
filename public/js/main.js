darkModeToggleBtn = document.querySelector("#invmode")
  darkModeToggleBtn.addEventListener("click",()=>{
  // console.log("clicked",document.getElementsByTagName('html')[0].style["filter"])

    document.getElementsByTagName('html')[0].classList.toggle("inverted")

// document.getElementsByTagName('html')[0].style["filter"] = "invert(1)"
// toggle value in the local storage
if(localStorage.getItem("darkMode") === 'true'){
  localStorage.setItem("darkMode",false)
}else{

  localStorage.setItem("darkMode",true)
  }



})

/* Dark mode */
if (
    (!localStorage["darkMode"] &&
        window.matchMedia("(prefers-color-scheme: dark)").matches) ||
   localStorage["darkMode"] === "true"
) {
    document.getElementsByTagName('html')[0].classList.toggle("inverted")
}
