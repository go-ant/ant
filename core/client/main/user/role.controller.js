angular
  .module('app')
  .controller('RoleListController', roleListCtrl)
  .controller('RoleController', roleCtrl);


/*
 * all roles
 */
roleListCtrl.$inject = ['$scope', '$mdDialog', '$translate', 'commService', 'roleService'];
function roleListCtrl($scope, $mdDialog, $translate, commService, roleService) {
  let vm = this;
  if (!$scope.isGranted('browse-roles')){
    return false;
  }
  $scope.initViewModel(vm);

  vm.remove = removeRole;
  vm.allowEditRoles = vm.isGranted('edit-roles');
  vm.allowDeleteRoles = vm.isGranted('delete-roles');
  getRoleList();

  function getRoleList() {
    return roleService.list()
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.roles = json.data;
          vm.pagination = json.pagination;
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  function removeRole(id, name) {
    var confirm = $mdDialog.confirm()
      .content($translate.instant('dlg.would_you_like_to_delete_role', { name: name}))
      .ok($translate.instant('delete'))
      .cancel($translate.instant('cancel'));
    $mdDialog.show(confirm).then(function() {
      roleService.remove(id)
        .then(function(json) {
          $scope.showToast(json.error.message, json.error);
          if (json.error.code == 0) {
            for(let i = 0; i < vm.roles.length; i++) {
              if (vm.roles[i].id == id) {
                vm.roles.splice(i, 1);
                break;
              }
            }
          }
        }, function(error) {commService.ajaxFailedCallBack($scope, error)});
    });
  }
}

/*
 * create/edit role
 */
roleCtrl.$inject = ['$scope', '$translate', 'commService', 'roleService'];
function roleCtrl($scope, $translate, commService, roleService) {
  let vm = this
    , listPermissionId = [];
  vm.roleId = $scope.$stateParams.roleId || 0;
  if (!$scope.isGranted(['add-roles', 'edit-roles'])){
    return false;
  }
  $scope.initViewModel(vm);

  vm.save = save;
  vm.existsPermission = existsPermission;
  vm.togglePermission = togglePermission;
  vm.allowSave = allowSave;
  vm.isSaving = false;
  vm.models = {};
  vm.models.role = {id: 0, name: ''};
  vm.models.permissions = {};

  if (vm.roleId == 0) {
    getPermissions();
  } else {
    getRole();
  }

  // get roles data
  function getRole() {
    roleService.get(vm.roleId, 'permission')
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0 && json.data) {
          vm.models.role = json.data;
          if (json.data.permissions) {
            listPermissionId = json.data.permissions.map(x => x.id);
          }
          getPermissions();
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  // get permissions data
  function getPermissions() {
    roleService.listPermission()
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          let permissions = {'core': [], 'additional': []};
          json.data.forEach(function(obj) {
            if (obj.is_core) {
              permissions.core.push(obj);
            } else {
              permissions.additional.push(obj);
            }
          });
          vm.models.permissions = permissions;
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  // check if id exists in permissions
  function existsPermission(permission) {
    return listPermissionId.indexOf(permission.id) > -1;
  }

  // toggle the permission option
  function togglePermission(permissionId) {
    var idx = listPermissionId.indexOf(permissionId);
    if (idx > -1) {
      listPermissionId.splice(idx, 1);
    }  else {
      listPermissionId.push(permissionId);
    }
  }

  function allowSave() {
    return (vm.roleId == 0 && vm.isGranted('add-roles')) || (vm.roleId > 0 && vm.isGranted('edit-roles'))
  }

  // save role
  function save() {
    if (vm.models.role.name.length == 0) {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.role_name_not_be_empty')});
      return false;
    }
    if (listPermissionId.length == 0) {
      $scope.showToast('', { code: -1, message: $translate.instant('msg.perm_not_be_empty')});
      return false;
    }

    vm.models.role.permissions = listPermissionId;
    vm.isSaving = true;
    if (vm.models.role.id == 0) {
      roleService.create(vm.models.role)
        .then(function(json) {
          vm.isSaving = false;
          $scope.showToast(json.error.message, json.error);
          if (json.error.code == 0) {
            vm.models.role = json.data;
            vm.breadcrumb = 'role_edit';
          }
        }, function(error) {vm.isSaving = false;commService.ajaxFailedCallBack($scope, error)});
    } else {
      roleService.edit(vm.models.role)
        .then(function(json) {
          vm.isSaving = false;
          $scope.showToast(json.error.message, json.error);
        }, function(error) {vm.isSaving = false;commService.ajaxFailedCallBack($scope, error)});
    }
  }

}
