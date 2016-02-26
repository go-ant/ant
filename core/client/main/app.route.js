angular
  .module('app')
  .config(appRouter);

appRouter.$inject = ['$stateProvider', '$urlRouterProvider', '$locationProvider', 'commServiceProvider'];
function appRouter ($stateProvider, $urlRouterProvider, $locationProvider, commServiceProvider) {
  let commService = commServiceProvider.$get();
  $urlRouterProvider.when('', '/posts');
  //$locationProvider.html5Mode(true);
  $stateProvider
    // posts
    .state('posts', {
      url: '/posts',
      templateUrl: commService.getTpls('posts/post_list'),
      controller: 'PostListController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'post_list'
      }
    })
    .state('posts/add', {
      url: '/posts/add',
      templateUrl: commService.getTpls('posts/post'),
      controller: 'PostController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'post_add'
      }
    })
    .state('posts/edit', {
      url: '/posts/{postId:[0-9]*}',
      templateUrl: commService.getTpls('posts/post'),
      controller: 'PostController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'post_edit'
      }
    })

    // users
    .state('profile', {
      url: '/profile',
      templateUrl: commService.getTpls('users/profile'),
      controller: 'UserEditController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: { title: 'nav.view_profile' }
    })
    .state('users', {
      url: '/users',
      templateUrl: commService.getTpls('users/user_list'),
      controller: 'UserListController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'user_list'
      }
    })
    .state('users/add', {
      url: '/users/add',
      templateUrl: commService.getTpls('users/user'),
      controller: 'UserCreateController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'user_add'
      }
    })
    .state('users/edit', {
      url: '/users/{userId:[0-9]*}',
      templateUrl: commService.getTpls('users/profile'),
      controller: 'UserEditController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: { title: 'user_edit' }
    })

    // roles
    .state('roles', {
      url: '/roles',
      templateUrl: commService.getTpls('users/role_list'),
      controller: 'RoleListController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'role_list'
      }
    })
    .state('roles/add', {
      url: '/roles/add',
      templateUrl: commService.getTpls('users/role'),
      controller: 'RoleController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'role_add'
      }
    })
    .state('roles/edit', {
      url: '/roles/{roleId:[0-9]*}',
      templateUrl: commService.getTpls('users/role'),
      controller: 'RoleController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'role_edit'
      }
    })

    // settings
    .state('settings/general', {
      url: '/general',
      templateUrl: commService.getTpls('settings/general'),
      controller: 'SettingGeneralController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'setting_general'
      }
    })
    .state('settings/navigation', {
      url: '/navigation',
      templateUrl: commService.getTpls('settings/navigation'),
      controller: 'SettingNavigationController',
      controllerAs: 'vm',
      resolve: {init: ['appService', function(appService) {return appService.promise;}]},
      data: {
        title: 'setting_navigation'
      }
    });
}


