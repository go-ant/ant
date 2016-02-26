angular
  .module('app')
  .controller('SettingNavigationController', settingNavigationCtrl);

settingNavigationCtrl.$inject = ['$scope', 'commService', 'settingService'];
function settingNavigationCtrl($scope, commService, settingService) {
  let vm = this;
  if (!$scope.isGranted('edit-settings')){
    return false;
  }
  $scope.initViewModel(vm);
  vm.save = save;
  vm.addNav = addNav;
  vm.removeNav = removeNav;
  vm.models = {};
  vm.models.setting = {};
  vm.models.nav = { label: '', url: ''};
  getSetting();

  function getSetting() {
    return settingService.list()
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.models.setting = json.data;
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  function addNav() {
    if (vm.models.nav.label != '' && vm.models.nav.url != '') {
      vm.models.setting.navigation.push(angular.copy(vm.models.nav));
      vm.models.nav.label = '';
      vm.models.nav.url = '';
      angular.element('#txt-label').focus();
    }
  }
  function removeNav(idx) {
    vm.models.setting.navigation.splice(idx, 1)
  }
  function save() {
    vm.isSaving = true;
    settingService.edit(vm.models.setting)
      .then(function(json) {
        vm.isSaving = false;
        $scope.showToast(json.error.message, json.error);
      }, function(error) {vm.isSaving=false;commService.ajaxFailedCallBack($scope, error)});
  }
}
