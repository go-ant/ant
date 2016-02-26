angular
  .module('app')
  .provider('commService', commService);

function commService() {
  let supportedLang = {
    'en': 'en - English',
    'zh-cn': 'zh-cn - 中文(简体)'
  };

  // custom dialog controller with cancel event
  let dialogCtrl = ['$mdDialog', function($mdDialog) {
    let vm = this;
    vm.cancel = cancel;

    function cancel() {
      $mdDialog.cancel();
    }
  }];

  // image upload dialog controller
  let uploadImageCtrl = function($scope, source, callback, cancelCallback) {

    return ['$mdDialog', function($mdDialog) {
      let vm = this;

      vm.success = success;
      vm.error = error;
      vm.toggleImgZoneMode = toggleImgZoneMode;
      vm.saveImageFromUrl = saveImageFromUrl;
      vm.deleteImage = deleteImage;
      vm.save = save;
      vm.cancel = cancel;
      vm.allowUpload = true;
      vm.imgFromUrl = false;
      vm.source = source || '';
      vm.sourceUrl = '';

      if (vm.source != '') {
        vm.allowUpload = false;
      }

      function success($file, message) {
        let json = JSON.parse(message);
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.source = json.data.file;
          vm.allowUpload = false;
        }
      }

      function error() {
        $scope.showToast('', { code: 10, message: 'The server returned an error (Server was not available)'});
      }

      function toggleImgZoneMode(e) {
        e.stopPropagation();
        vm.imgFromUrl = !vm.imgFromUrl;
      }

      function saveImageFromUrl() {
        if (vm.sourceUrl != '') {
          vm.source = vm.sourceUrl;
          vm.allowUpload = false;
        }
        vm.imgFromUrl = false;
      }

      function deleteImage() {
        vm.source = '';
        vm.allowUpload = true;
      }

      function save() {
        callback && callback(vm.source);
        $mdDialog.cancel();
      }

      function cancel() {
        cancelCallback && cancelCallback();
        $mdDialog.cancel();
      }
    }];
  };


  this.$get = function() {
    return {
      supportedLang: supportedLang,
      getSysLang: getSysLang,
      ajaxContentCallBack: ajaxContentCallBack,
      ajaxFailedCallBack: ajaxFailedCallBack,
      getTpls: getTpls,
      dialogCtrl: dialogCtrl,
      uploadImageCtrl: uploadImageCtrl
    }
  };

  function ajaxContentCallBack(res) {
    return res.data;
  }

  function ajaxFailedCallBack($scope) {
    $scope.showToast('', { code: 10, message: 'The server returned an error (Server was not available)'});
  }

  function getTpls(tpl) {
    return tpl + '.html';
  }

  function getSysLang() {
    var lang = navigator.languages? navigator.languages[0] : (navigator.language || navigator.userLanguage);
    lang = lang.toLowerCase();
    if (!supportedLang[lang]) {
      lang = 'en';
    }
    return lang;
  }

}
