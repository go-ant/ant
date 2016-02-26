angular
  .module('app')
  .factory('userService', userService);

userService.$inject = ['$http', 'commService'];
function userService($http, commService) {
  let service = {
    setup: setup,
    login: login,
    list: list,
    get: get,
    create: create,
    edit: edit,
    changePassword: changePassword
  };
  return service;


  function setup(app) {
    return $http.post('../api/setup', app)
      .then(commService.ajaxContentCallBack);
  }

  function login(user) {
    return $http.post('../api/signin', {name: user.name, password: user.password})
      .then(commService.ajaxContentCallBack);
  }

  function list(opts) {
    return $http.get('../api/backend/users', {params: {limit: opts.limit || 15, page: opts.page || 1, include: opts.include}})
      .then(commService.ajaxContentCallBack);
  }

  function get(id, include) {
    return $http.get('../api/backend/users/info/' + id, {params: {include: include}})
      .then(commService.ajaxContentCallBack);
  }

  function create(user) {
    return $http.post('../api/backend/users/info', {name: user.name, password: user.password, role_id: user.roles[0].id})
      .then(commService.ajaxContentCallBack);
  }

  function edit(user) {
    return $http.put('../api/backend/users/info/' + user.id, {
        name: user.name,
        slug: user.slug,
        email: user.email,
        location: user.location,
        website: user.website,
        bio: user.bio,
        avatar: user.avatar,
        cover: user.cover,
        role_id: user.roles[0].id
      })
      .then(commService.ajaxContentCallBack);
  }

  function changePassword(user_id, old_passwd, new_passwd, verify_passwd) {
    return $http.put('../api/backend/users/password/' + user_id, { old_password: old_passwd, new_password: new_passwd, verify_password: verify_passwd })
      .then(commService.ajaxContentCallBack);

  }
}