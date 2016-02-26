angular
  .module('app')
  .controller('SettingGeneralController', settingGeneralCtrl);

settingGeneralCtrl.$inject = ['$scope', '$mdDialog', 'commService', 'themeService', 'settingService'];
function settingGeneralCtrl($scope, $mdDialog, commService, themeService, settingService) {
  let vm = this;
  if (!$scope.isGranted('edit-settings')){
    return false;
  }
  $scope.initViewModel(vm);

  vm.save = save;
  vm.showImageUpload = showImageUpload;
  vm.models = {};
  vm.models.setting = {
    theme: 'journey'
  };
  vm.models.themes = [];
  vm.models.languages = commService.supportedLang;

  getThemeList();
  getSetting();


  function showImageUpload(ev, type) {
    let source = '';
    if (type == 'logo') {
      if (vm.models.setting.logo != '') {
        source = vm.models.setting.logo;
      }
    } else {
      if (vm.models.setting.cover != '') {
        source = vm.models.setting.cover
      }
    }
    $mdDialog.show({
      controller: commService.uploadImageCtrl($scope, source, saveImg(type)),
      controllerAs: 'vm',
      templateUrl: commService.getTpls('tmpl/image_upload'),
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose:true
    });
  }

  function saveImg(type) {
    return function(url) {
      if (type == 'logo') {
        vm.models.setting.logo = url;
      } else {
        vm.models.setting.cover = url;
      }
    }
  }

  function getThemeList() {
    return themeService.list()
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.models.themes = json.data;
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  function getSetting() {
    return settingService.list()
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.models.setting = json.data;
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }


  function save() {
    vm.isSaving = true;
    settingService.edit(vm.models.setting)
      .then(function(json) {
        vm.isSaving = false;
        $scope.showToast(json.error.message, json.error);
        if (json.error.code == 0) {
          $scope.setLanguage(vm.models.setting.language);
        }
      }, function(error) {vm.isSaving=false;commService.ajaxFailedCallBack($scope, error)});
  }
}
