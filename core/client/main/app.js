 angular.module('app', [
    'ngSanitize',
    'ngMaterial',
    'ngMessages',
    'ui.router',
    'LocalStorageModule',
    'angular-loading-bar',
    'flow',
    'pascalprecht.translate'
  ]).run(appRun);

/*
 * base configuration
 */
appRun.$inject = ['$rootScope', '$mdMedia', '$mdSidenav', '$mdToast', '$state', '$stateParams', '$translate'];
function appRun($rootScope, $mdMedia, $mdSidenav, $mdToast, $state, $stateParams, $translate) {
  let toastDelay = 5000
    , userInfo = {};
  $rootScope.$state = $state;
  $rootScope.$stateParams = $stateParams;
  $rootScope.showToast = showToast;
  $rootScope.initViewModel = initViewModel;
  $rootScope.isGranted = isGranted;
  $rootScope.setUserInfo = setUserInfo;
  $rootScope.setLanguage = setLanguage;

  // responsive toggle menus
  function openMenu() {
    if (!$mdMedia('gt-md')) {
      $mdSidenav('menu').toggle();
    }
  }

  $rootScope.$on('$stateChangeStart', function(event, toState, toParams, fromState, fromParams){
    if (!fromState.abstract && !$mdMedia('gt-md')) {
      $mdSidenav('menu').close();
    }
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

  function setLanguage(lang) {
    $translate.use(lang);
  }

  // set login user info and permissions
  function setUserInfo(user) {
    userInfo = user;
  }

  function initViewModel(vm) {
    vm.openMenu = openMenu;
    vm.breadcrumb = $state.current.data.title;
    vm.isGranted = checkPermission;
    vm.userInfo = userInfo;
  }

  function isGranted(perm) {
    if (!checkPermission(perm)) {
      $state.go('dashboard');
      showToast('', {message: $translate.instant('msg.no_permission')});
      return false;
    }
    return true;
  }

  function checkPermission(perm) {
    if (angular.isArray(perm)) {
      for(let i = 0; i < userInfo.permissions.length; i++) {
        if (perm.indexOf(userInfo.permissions[i]) >= 0) {
          return true
        }
      }
    }
    return userInfo.permissions.indexOf(perm) >= 0;
  }


}
