angular.module('app', [
    'ngMaterial',
    'flow',
    'pascalprecht.translate'
  ])
  .controller('SetupController', setupCtrl);


setupCtrl.$inject = ['$scope', '$mdToast', '$translate', 'commService', 'userService'];
function setupCtrl($scope, $mdToast, $translate, commService, userService) {
  let vm = this
   , toastDelay = 5000;

  $scope.showToast = showToast;
  vm.changeLang = changeLang;
  vm.setup = setup;
  vm.isSaving = false;
  vm.models = {};
  vm.models.languages = commService.supportedLang;
  vm.models.app = {
    name: '',
    password: '',
    title: '',
    avatar: '',
    language: commService.getSysLang()
  };

  $translate.use(vm.models.app.language);

  $scope.$on('flow::fileAdded', function (event, $flow, flowFile) {
   var options = {
     maxWidth: 200,
     maxHeight: 200,
     canvas: true,
     noRevoke: true
   };
   loadImage.parseMetaData(flowFile.file, function (data) {
     if (data.exif) {
       options.orientation = data.exif.get('Orientation');
     }
     loadImage(
       flowFile.file,
       function (canvas) {
         vm.models.app.avatar = canvas.toDataURL();
         $scope.$apply();
       }, options);
   });
  });

  function showToast(msgSuccessful, error, showSucc) {
    switch (error.code) {
      case 0:
        if (showSucc == undefined || showSucc) {
          $mdToast.show($mdToast.simple().theme('success').position('bottom right').content($translate.instant(msgSuccessful)).hideDelay(toastDelay));
        }
        break;
      case 10:
        $mdToast.show($mdToast.simple().theme('fail').position('bottom right').content($translate.instant(error.message)).hideDelay(toastDelay));
        break;
      default:
        $mdToast.show($mdToast.simple().theme('warn').position('bottom right').content($translate.instant(error.message)).hideDelay(toastDelay));
        break;
    }
  }

  function changeLang() {
    $translate.use(vm.models.app.language);
  }

  function setup() {
    if (vm.models.app.name.length == 0) {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.user_name_not_be_empty')});
      return false;
    }
    if (/^\d+$/.test(vm.models.app.name)) {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.not_support_numeric_user')});
      return false;
    }
    if (vm.models.app.password.length == 0) {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.user_password_not_be_empty')});
      return false;
    }
    if (vm.models.app.password.length < 8) {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.password_at_least_8_char')});
      return false;
    }

    vm.isSaving = true;
    userService.setup(vm.models.app)
      .then(function(json) {
        vm.isSaving = false;
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          location.href = './';
        }
      }, function(error) {vm.isSaving=false;commService.ajaxFailedCallBack($scope, error)});
  }
}
