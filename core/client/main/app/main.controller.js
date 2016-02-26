angular
  .module('app')
  .controller('MainController', mainCtrl);

mainCtrl.$inject = ['$scope', 'appService'];
function mainCtrl($scope, appService) {
  let vm = this;
  vm.models = {};
  // profile menu
  vm.activeProfileMenu = false;
  vm.toggleProfileMenu = toggleProfileMenu;

  initApp();

  function toggleProfileMenu() {
    vm.activeProfileMenu = !vm.activeProfileMenu;
  }

  function initApp() {
    appService.init()
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.models.menus = json.data.menus;
          vm.models.user = json.data.user;

          let user = angular.copy(json.data.user);
          user.permissions = json.data.permissions;
          $scope.setUserInfo(user);

          // init interface language
          $scope.setLanguage(json.data.language);

          appService.defer.resolve();
        }
      }, function(error) {vm.isSaving=false;commService.ajaxFailedCallBack($scope, error)});
  }
}
