'use strict'

const { app, Menu, BrowserWindow, dialog } = require('electron');

var log = require('electron-log');

const path = require('path');

const childProcess = require('child_process');

const cwd = require('process').cwd();

// This adds refresh and devtools console keybindings
// Page can refresh with cmd+r, ctrl+r, F5
// Devtools can be toggled with cmd+alt+i, ctrl+shift+i, F12
require('electron-debug')({enabled: true, showDevTools: false});
require('electron-context-menu')({});


global.eval = function() { throw new Error('bad!!'); }

const defaultURL = 'http://127.0.0.1:8642/';
let currentURL;

// Force everything localhost, in case of a leak
app.commandLine.appendSwitch('host-rules', 'MAP * 127.0.0.1, EXCLUDE api.coinmarketcap.com, api.github.com');
app.commandLine.appendSwitch('ssl-version-fallback-min', 'tls1.2');
app.commandLine.appendSwitch('--no-proxy-server');
app.setAsDefaultProtocolClient('unicoin');



// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win;

var unicoin = null;

function startUnicoin() {
  console.log('Starting unicoin from electron');

  if (unicoin) {
    console.log('Unicoin already running');
    app.emit('unicoin-ready');
    return
  }

  var reset = () => {
    unicoin = null;
  }

  // Resolve unicoin binary location
  var appPath = app.getPath('exe');
  var exe = (() => {
        switch (process.platform) {
  case 'darwin':
    return path.join(appPath, '../../Resources/app/unicoin');
  case 'win32':
    // Use only the relative path on windows due to short path length
    // limits
    return './resources/app/unicoin.exe';
  case 'linux':
    return path.join(path.dirname(appPath), './resources/app/unicoin');
  default:
    return './resources/app/unicoin';
  }
})()

  var args = [
    '-launch-browser=false',
    '-gui-dir=' + path.dirname(exe),
    '-color-log=false', // must be disabled for web interface detection
    '-logtofile=true',
    '-download-peerlist=true',
    '-enable-seed-api=true',
    '-enable-wallet-api=true',
    '-rpc-interface=false',
    "-disable-csrf=false"
    // will break
    // broken (automatically generated certs do not work):
    // '-web-interface-https=true',
  ]
  unicoin = childProcess.spawn(exe, args);

  unicoin.on('error', (e) => {
    dialog.showErrorBox('Failed to start unicoin', e.toString());
    app.quit();
  });

  unicoin.stdout.on('data', (data) => {
    console.log(data.toString());
    // Scan for the web URL string
    if (currentURL) {
      return
    }
    const marker = 'Starting web interface on ';
    var i = data.indexOf(marker);
    if (i === -1) {
      return
    }
    currentURL = defaultURL;
    app.emit('unicoin-ready', { url: currentURL });
  });

  unicoin.stderr.on('data', (data) => {
    console.log(data.toString());
  });

  unicoin.on('close', (code) => {
    // log.info('Unicoin closed');
    console.log('Unicoin closed');
    reset();
  });

  unicoin.on('exit', (code) => {
    // log.info('Unicoin exited');
    console.log('Unicoin exited');
    reset();
  });
}

function createWindow(url) {
  if (!url) {
    url = defaultURL;
  }

  // To fix appImage doesn't show icon in dock issue.
  var appPath = app.getPath('exe');
  var iconPath = (() => {
    switch (process.platform) {
      case 'linux':
        return path.join(path.dirname(appPath), './resources/icon512x512.png');
    }
  })()

  // Create the browser window.
  win = new BrowserWindow({
    width: 1200,
    height: 900,
    title: 'Unicoin',
    icon: iconPath,
    nodeIntegration: false,
    webPreferences: {
      webgl: false,
      webaudio: false,
    },
  });

  // patch out eval
  win.eval = global.eval;

  const ses = win.webContents.session
  ses.clearCache(function () {
    console.log('Cleared the caching of the unicoin wallet.');
  });

  ses.clearStorageData([],function(){
    console.log('Cleared the stored cached data');
  });

  win.loadURL(url);

  // Open the DevTools.
  // win.webContents.openDevTools();

  // Emitted when the window is closed.
  win.on('closed', () => {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    win = null;
  });

  win.webContents.on('will-navigate', function(e, url) {
    e.preventDefault();
    require('electron').shell.openExternal(url);
  });

  // create application's main menu
  var template = [{
    label: "Unicoin",
    submenu: [
      { label: "About Unicoin", selector: "orderFrontStandardAboutPanel:" },
      { type: "separator" },
      { label: "Quit", accelerator: "Command+Q", click: function() { app.quit(); } }
    ]
  }, {
    label: "Edit",
    submenu: [
      { label: "Undo", accelerator: "CmdOrCtrl+Z", selector: "undo:" },
      { label: "Redo", accelerator: "Shift+CmdOrCtrl+Z", selector: "redo:" },
      { type: "separator" },
      { label: "Cut", accelerator: "CmdOrCtrl+X", selector: "cut:" },
      { label: "Copy", accelerator: "CmdOrCtrl+C", selector: "copy:" },
      { label: "Paste", accelerator: "CmdOrCtrl+V", selector: "paste:" },
      { label: "Select All", accelerator: "CmdOrCtrl+A", selector: "selectAll:" }
    ]
  }];

  Menu.setApplicationMenu(Menu.buildFromTemplate(template));
}

// Enforce single instance
const alreadyRunning = app.makeSingleInstance((commandLine, workingDirectory) => {
      // Someone tried to run a second instance, we should focus our window.
      if (win) {
        if (win.isMinimized()) {
          win.restore();
        }
        win.focus();
      } else {
        createWindow(currentURL || defaultURL);
}
});

if (alreadyRunning) {
  app.quit();
  return;
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', startUnicoin);

app.on('unicoin-ready', (e) => {
  createWindow(e.url);
});

// Quit when all windows are closed.
app.on('window-all-closed', () => {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
  app.quit();
}
});

app.on('activate', () => {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (win === null) {
  createWindow();
}
});

app.on('will-quit', () => {
  if (unicoin) {
    unicoin.kill('SIGINT');
  }
});

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
