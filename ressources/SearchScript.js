document.getElementById("dateCheckBox").addEventListener("change", function () {
    let dateStart = document.getElementById("date1");
    let dateEnd = document.getElementById("date2");
    if (this.checked) {
        content.classList.remove("disabled");
    } else {
        content.classList.add("disabled");
    }
});