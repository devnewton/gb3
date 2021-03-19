let linuxfrInfo = document.getElementById('linuxfr-info');
try {
    let tokenStr = window.atob(window.location.hash.substr(1));
    let token = JSON.parse(tokenStr);
    localStorage.linuxfr_access_token = token.access_token;
    localStorage.linuxfr_refresh_token = token.refresh_token;
    localStorage.linuxfr_expires_in = token.expires_in;
    linuxfrInfo.innerText = "Authorization token successfullty received from linuxfr :-)";
  }
  catch(error) {
    console.error(error);
    linuxfrInfo.innerText = "Cannot read authorization token from linuxfr :-(";
  }