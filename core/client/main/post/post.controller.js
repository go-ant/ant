angular
  .module('app')
  .controller('PostListController', postListCtrl)
  .controller('PostController', postCtrl);

/*
 * all posts
 */
postListCtrl.$inject = ['$scope', '$mdDialog', '$translate', 'postService'];
function postListCtrl($scope, $mdDialog, $translate, postService) {
  let vm = this;
  if (!$scope.isGranted(['browse-posts', 'edit-all-posts'])){
    return false;
  }
  $scope.initViewModel(vm);
  vm.getPostList = getPostList;
  vm.remove = removePost;
  vm.allowEditPosts = vm.isGranted(['edit-posts', 'edit-all-posts']);
  vm.models = {};
  getPostList();

  function getPostList(flag, page) {
    return postService.list({page: page || 1})
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.models.posts = json.data;
          vm.pagination = json.pagination;
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  function removePost(id, title) {
    var confirm = $mdDialog.confirm()
      .htmlContent($translate.instant('dlg.would_you_like_to_delete_post', { title: title}))
      .ok($translate.instant('delete'))
      .cancel($translate.instant('cancel'));
    $mdDialog.show(confirm).then(function() {
      postService.remove(id)
        .then(function(json) {
          $scope.showToast(json.error.message, json.error);
          if (json.error.code == 0) {
            for(let i = 0; i < vm.models.posts.length; i++) {
              if (vm.models.posts[i].id == id) {
                vm.models.posts.splice(i, 1);
                break;
              }
            }
          }
        }, function(error) {commService.ajaxFailedCallBack($scope, error)});
    });
  }
}


/*
 * create/edit posts
 */
postCtrl.$inject = ['$scope', '$timeout', '$mdSidenav', '$mdDialog', '$translate', 'commService', 'postService'];
function postCtrl($scope, $timeout, $mdSidenav, $mdDialog, $translate, commService, postService) {
  let vm = this
    , converter
    , uploadMaxSize = 5<<20
    , $coverZone = $('.post-cover-zone')
    , btnText = {
        'published': {'published': 'update_post', 'draft': 'unpublish'},
        'draft': {'published': 'publish_now', 'draft': 'save_draft'}
      };

  vm.postId = $scope.$stateParams.postId || 0;
  if (!$scope.isGranted(['add-posts', 'edit-posts', 'edit-all-posts'])) {
    return false;
  }
  $scope.initViewModel(vm);

  vm.remove = removePost;
  vm.save = save;
  vm.showMarkdownHelpe = showMarkdownHelpe;
  vm.toggleSettings = toggleSettings;
  vm.toggleView = toggleView;
  vm.toggleImgZoneMode = toggleImgZoneMode;
  vm.setStatus = setStatus;
  vm.saveCoverFromUrl = saveCoverFromUrl;
  vm.cancelPostCover = cancelPostCover;
  vm.btnText = 'save_draft';
  vm.postStatus = 'draft';
  vm.postCover = '';
  vm.postTags = [];
  vm.showMarkdown = true;
  vm.isSaving = false;
  vm.models = {};
  vm.models.post = {id: 0, status: 'draft'};

  if (vm.postId > 0) {
    getPost();
  }

  converter = new showdown.Converter({
    tables: true,
    extensions: ['github', 'antimagepreview']
  });

  initUploadPostCover();

  $scope.$watch('vm.models.post.markdown', function(){
    vm.models.post.html = converter.makeHtml(vm.models.post.markdown);
    $timeout(function() {
      let imgs = $('.post-img-zone');
      if (imgs.length > 0) {
        uploadPostImage(imgs)
      }
    }, 20);

  });


  function showMarkdownHelpe(ev) {
    $mdDialog.show({
      controller: commService.dialogCtrl,
      controllerAs: 'vm',
      templateUrl: commService.getTpls('tmpl/markdown_help'),
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose:true
    });
  }

  function toggleSettings() {
    $mdSidenav('settings').toggle();
  }

  function toggleView() {
    vm.showMarkdown = !vm.showMarkdown
  }

  function toggleImgZoneMode(e) {
    e.stopPropagation();
    let $zone = $(e.target).closest('.img-zone');
    if (!$zone.hasClass('img-source-url')) {
      $zone.addClass('img-source-url');
    } else {
      $zone.removeClass('img-source-url');
    }
  }

  function initFlow() {
    return new Flow({
      target: '../api/backend/upload',
      chunkSize: uploadMaxSize,
      testChunks: false,
      singleFile: true
    });
  }

  function uploadPostImage(imgs) {
    imgs.each(function(idx) {
      let elementIndex = idx
        , flow = initFlow();
      flow.on('fileAdded', function (file){
        if (file.size > uploadMaxSize){
          $scope.showToast('', { message: 'file_too_large'});
          return false;
        }
      });
      flow.on('filesSubmitted', function() {
        imgs.eq(idx).addClass('img-uploading');
        flow.upload();
      });
      flow.on('fileSuccess', function(file, message){
        let json = JSON.parse(message);
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          let nth = 0
            , newMarkdown = vm.models.post.markdown.replace(/^(?:\{<(.*?)>\})?!(?:\[([^\n\]]*)\])(:\(([^\n\]]*)\))?$/gim, function (match){
            nth++;
            return (nth === (elementIndex+1)) ? `${match}(${json.data.file})` : match;
          });
          vm.models.post.markdown = newMarkdown;
          $scope.$apply();
        }
      });
      flow.assignBrowse($(this).find('.img-desc')[0], false, true, {accept: 'image/*'});
      flow.assignDrop($(this).find('.img-desc')[0]);
    });

  }

  function initUploadPostCover() {
    let flow = initFlow();
    flow.on('fileAdded', function (file){
      if (file.size > uploadMaxSize){
        $scope.showToast('', { message: 'file_too_large'});
        return false;
      }
    });
    flow.on('filesSubmitted', function() {
      $coverZone.addClass('img-uploading');
      flow.upload();
    });
    flow.on('fileSuccess', function(file, json){
      json = JSON.parse(json);
      $scope.showToast(json.error.message, json.error, false);
      $coverZone.removeClass('img-uploading');
      if (json.error.code == 0) {
        vm.models.post.cover = json.data.file;
        vm.postCover = json.data.file;
        $coverZone.addClass('img-uploaded');
        $scope.$apply();
      }
    });
    flow.assignBrowse($coverZone.find('.img-desc')[0], false, true, {accept: 'image/*'});
    flow.assignDrop($coverZone.find('.img-desc')[0]);
  }

  function saveCoverFromUrl() {
    vm.postCover = vm.models.post.cover;
    $coverZone.removeClass('img-source-url');
    $coverZone.addClass('img-uploaded');
  }

  function cancelPostCover() {
    vm.postCover = '';
    vm.models.post.cover = '';
    $coverZone.removeClass('img-uploaded');
  }

  function setStatus(status) {
    vm.postStatus = status;
    vm.btnText = btnText[vm.models.post.status][status];
  }

  function getPost() {
    return postService.get(vm.postId)
      .then(function(json) {
        $scope.showToast(json.error.message, json.error, false);
        if (json.error.code == 0) {
          vm.models.post = json.data;
          vm.postCover = vm.models.post.cover;
          setStatus(vm.models.post.status);
        }
      }, function(error) {commService.ajaxFailedCallBack($scope, error)});
  }

  function removePost(title) {
    var confirm = $mdDialog.confirm()
      .htmlContent($translate.instant('dlg.would_you_like_to_delete_post', { title: title}))
      .ok($translate.instant('delete'))
      .cancel($translate.instant('cancel'));
    $mdDialog.show(confirm).then(function() {
      postService.remove(vm.models.post.id)
        .then(function(json) {
          $scope.showToast(json.error.message, json.error);
          if (json.error.code == 0) {
            $scope.$state.go('posts');
          }
        }, function(error) {commService.ajaxFailedCallBack($scope, error)});
    });
  }

  // save post
  function save() {
    vm.isSaving = true;
    vm.models.post.status = vm.postStatus;
    vm.btnText = btnText[vm.models.post.status][vm.postStatus];

    if (vm.models.post.id == 0) {
      postService.create(vm.models.post)
        .then(function(json) {
          vm.isSaving = false;
          $scope.showToast(json.error.message, json.error);
          if (json.error.code == 0) {
            vm.models.post = json.data;
          }
        }, function(error) {vm.isSaving = false;commService.ajaxFailedCallBack($scope, error)});
    } else {
      postService.edit(vm.models.post)
        .then(function(json) {
          vm.isSaving = false;
          $scope.showToast(json.error.message, json.error);
          if (json.error.code == 0) {
            vm.models.post = json.data;
          }
        }, function(error) {vm.isSaving = false;commService.ajaxFailedCallBack($scope, error)});
    }
  }

}
