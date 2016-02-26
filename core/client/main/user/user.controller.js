angular
  .module('app')
  .controller('UserListController', userListCtrl)
  .controller('UserCreateController', userCtrl)
  .controller('UserEditController', userEditCtrl);

/*
 * all users
 */
userListCtrl.$inject = ['$scope', 'commService', 'userService'];
function userListCtrl($scope, commService, userService) {
  let vm = this;
  if (!$scope.isGranted(['browse-users', 'browse-all-users'])){
    return false;
  }
  $scope.initViewModel(vm);

  getUserList();

  function getUserList(flag, page) {
    return userService.list({limit: 15, page: page, include: 'role'})
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.users = json.data;
          vm.pagination = json.pagination;
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }
}

/*
 * create user
 */
userCtrl.$inject = ['$scope', '$translate', 'commService', 'roleService', 'userService'];
function userCtrl($scope, $translate, commService, roleService, userService) {
  let vm = this;
  if (!$scope.isGranted('add-users')){
    return false;
  }
  $scope.initViewModel(vm);

  vm.save = save;
  vm.getRoles = getRoles;
  vm.isSaving = false;
  vm.models = {};
  vm.models.user = {id: 0, name: '', password: '', roles:[]};

  // get roles data
  function getRoles() {
    roleService.list()
      .then(function(json) {
        vm.roles = json.data;
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  // save user
  function save() {
    if (vm.models.user.roles.length == 0) {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.role_not_be_empty')});
      return false;
    }
    if (vm.models.user.name == '') {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.user_name_not_be_empty')});
      return false;
    }
    if (/^\d+$/.test(vm.models.user.name)) {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.not_support_numeric_user')});
      return false;
    }
    if (vm.models.user.password == '') {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.user_password_not_be_empty')});
      return false;
    }

    vm.isSaving = true;
    userService.create(vm.models.user)
      .then(function(json) {
        vm.isSaving = false;
        $scope.showToast(json.error.message, json.error);
        if (json.error.code == 0) {
          $scope.$state.go('users/edit', {userId: json.data.id});
        }
      }, function(error) {vm.isSaving=false;commService.ajaxFailedCallBack($scope, error)});
  }

}

/*
 * edit user
 */
userEditCtrl.$inject = ['$scope', '$mdDialog', 'commService', 'roleService', 'userService'];
function userEditCtrl($scope, $mdDialog, commService, roleService, userService) {
  let vm = this
    , defAvatar = './assets/css/images/avatar.png'
    , defCover = './assets/css/images/cover.jpg';
  vm.userId = $scope.$stateParams.userId || 0;
  if (vm.userId > 0 && !$scope.isGranted('edit-users') && $scope.isGranted(vm)){
    return false;
  }
  $scope.initViewModel(vm);

  vm.save = save;
  vm.savePassword = savePassword;
  vm.showImageUpload = showImageUpload;
  vm.allowUpdateRole = allowUpdateRole;
  vm.allowUpdatePassword = allowUpdatePassword;
  vm.needPassword = needPassword;
  vm.allowSave = allowSave;
  vm.isSaving = false;
  vm.models = {};
  vm.password = {old: '', new: ''};

  getRoles();
  getUser();

  function showImageUpload(ev, type) {
    let source = '';
    if (type == 'avatar') {
      if (vm.models.user.avatar != defAvatar) {
        source = vm.models.user.avatar;
      }
    } else {
      if (vm.models.user.cover != defCover) {
        source = vm.models.user.cover
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
      if (type == 'avatar') {
        if (url == '') {
          vm.models.user.avatar = defAvatar;
        } else {
          vm.models.user.avatar = url;
        }
      } else {
        if (url == '') {
          vm.models.user.cover = defCover;
        } else {
          vm.models.user.cover = url;
        }
      }
    }
  }

  function getUser() {
    return userService.get(vm.userId, 'role')
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.models.user = json.data;
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  // get roles data
  function getRoles() {
    roleService.list()
      .then(function(json) {
        vm.roles = json.data;
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  function allowUpdateRole() {
    return vm.models.user.id != vm.userInfo.id && vm.models.user.roles[0].slug != 'owner' && vm.isGranted('assign-roles');
  }

  function allowUpdatePassword() {
    return vm.models.user.roles[0].slug != 'owner' || vm.userInfo.roles[0].slug == 'owner';
  }

  function needPassword() {
    return vm.models.user.id == vm.userInfo.id;
  }

  function allowSave() {
    return (vm.models.user.roles[0].slug == 'owner' && vm.userInfo.roles[0].slug == 'owner') ||
      vm.models.user.roles[0].slug != 'owner' && (vm.userId == 0 || vm.isGranted('edit-users'));
  }

  function save() {
    vm.isSaving = true;
    return userService.edit(vm.models.user)
      .then(function(json) {
        vm.isSaving = false;
        $scope.showToast(json.error.message, json.error);
      }, function(error) {vm.isSaving=false;commService.ajaxFailedCallBack($scope, error)});
  }

  function savePassword() {
    vm.isSaving = true;
    return userService.changePassword(vm.models.user.id, vm.password.old, vm.password.new, vm.password.verify)
      .then(function(json) {
        vm.isSaving = false;
        $scope.showToast(json.error.message, json.error);
        if (json.error.code == 0) {
          vm.password = {};
        }
      }, function(error) {vm.isSaving=false;commService.ajaxFailedCallBack($scope, error)});
  }

}

