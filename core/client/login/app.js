angular
  .module('app', [
    'ngMaterial',
    'pascalprecht.translate'
  ])
  .controller('LoginController', loginCtrl);


loginCtrl.$inject = ['$scope', '$mdToast', '$translate', 'commService', 'userService'];
 function loginCtrl($scope, $mdToast, $translate, commService, userService) {
   let vm = this
     , toastDelay = 5000;

   $scope.showToast = showToast;

   vm.login = login;
   vm.isSaving = false;
   vm.models = {};
   vm.models.user = {
     name: '',
     password: ''
   };

   $translate.use(commService.getSysLang());

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

   function login() {
     if (vm.models.user.name.length == 0 || vm.models.user.name.password == 0) {
       $scope.showToast('', { code: -1, message: $translate.instant('msg.fill_out_the_form_to_sign_in')});
       return false;
     }

     vm.isSaving = true;
     return userService.login(vm.models.user)
       .then(function(json) {
         vm.isSaving = false;
         $scope.showToast(json.error.message, json.error, false);
         if (json.error.code == 0) {
           location = './';
         }
       }, function(error) {vm.isSaving=false;commService.ajaxFailedCallBack($scope, error)});
   }
 }
