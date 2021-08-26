window.addEventListener('DOMContentLoaded', () => {
  let btn = document.getElementById("submit");
  btn.addEventListener("click", loadingMessage);

  let uploadFile = document.getElementById("upload");
  uploadFile.addEventListener("change", renderPreview);
});

function loadingMessage() {
  let txt = btn.innerText;
  txt = 'Please wait while your image is being processed...\n(this might take a while, depending on the resolution of your image)';
  btn.setAttribute("aria-busy", "true");
  btn.innerText = txt;
}

function renderPreview() {
  let file = uploadFile.files;
  if (file.length > 0) {
    let imageFile = file[0];
    checkFileType(imageFile);
    checkFileSize(imageFile);
    renderImg(imageFilee);
  } 
}

function checkFileType(file) {
    const fileType = file.type;
    if (fileType != "image/jpeg" && fileType != "image/png") {
      notie.alert({type: 'error', text:'The selected file does not seem to be an image file. Please try another file.'})
      return btn.disabled = true;
      } else {
          btn.disabled = false;
    }
}

function checkFileSize(file) {
    const fileSize = file.size;
    const fileSizeMB = Math.round((fileSize/1024/1024));
      if (fileSizeMB > 4) {
      notie.alert({type: 'error', text:'The selected file exceeds 4mb in size. Please select a smaller file.'})
        return btn.disabled = true;
      } else {
          btn.disabled = false;
      }
}

function renderImg(file) {
    var reader = new FileReader();
        reader.onload = () => {
        var img = document.createElement("img");
            img.onload = () => {
              var MAX_WIDTH = 640;
              var width = img.width;
              var height = img.height;

              if (width > MAX_WIDTH) {
                height = height * (MAX_WIDTH/width);
                width = MAX_WIDTH;
              }

              // Dynamically create a canvas element
              var canvas = document.createElement("canvas");
              canvas.width = width;
              canvas.height = height;

              // var canvas = document.getElementById("canvas");
              var ctx = canvas.getContext("2d");

              // Actual resizing
              ctx.drawImage(img, 0, 0, width, height);

              // Show resized image in preview element
              var dataurl = canvas.toDataURL(file.type);
              document.getElementById("preview").src = dataurl;
            }
        img.src = e.target.result;
      }
    reader.readAsDataURL(file);
}
