angular
  .module('app')
  .directive('menuLink', directiveMenuLink)
  .directive('menuToggle', directiveMenuToggle);


function directiveMenuLink() {
  return {
    restrict: 'E',
    scope: {
      section: '='
    },
    template: '<md-button ng-click="goToPage()" ui-sref="{{section.URL}}" aria-label="{{::section.Label}}"><div class="md-button-inner"><md-icon ng-if="::section.Icon" ng-bind="::section.Icon"></md-icon><span class="md-inline-list-icon-label" translate="{{section.Name}}"></span></div></md-button>',
    link: function($scope, $element, $attr) {
      $scope.goToPage = function() {
        $element.closest('menu-toggle').parent().find('menu-link.active').removeClass('active');
        $element.addClass('active');
      }
    }
  };
}

function directiveMenuToggle() {
  return {
    restrict: 'E',
    scope: {
      section: '='
    },
    template: '<md-button ng-if="::section.URL" ui-sref="{{section.URL}}" aria-label="{{::section.Label}}" ng-click="toggleMenu()"><div class="md-button-inner"><md-icon ng-if="::section.Icon" ng-bind="::section.Icon"></md-icon><span class="md-inline-list-icon-label" translate="{{section.Name}}"></span></div></md-button>' +
    '<md-button ng-if="::!section.URL" aria-label="{{::section.Label}}" ng-click="toggleMenu()"><div class="md-button-inner"><md-icon ng-if="::section.Icon" ng-bind="::section.Icon"></md-icon><span class="md-inline-list-icon-label" translate="{{section.Name}}"></span></div></md-button>' +
    '<md-list ng-if="::section.List" class="menu-toggle-list">' +
    '<md-list-item ng-repeat="menu in ::section.List"><menu-link section="menu" ng-class="{active: ($state.current.name == menu.URL)}"></menu-link></md-list-item>' +
    '</md-list>',
    controller: ['$scope', '$state', function($scope, $state) {
      $scope.$state = $state
    }],
    link: function($scope, $element) {

      $scope.toggleMenu = function() {
        if ($element.hasClass('active')) {
          $element.removeClass('active');
        } else {
          $element.parent().find('>.active').removeClass('active');
          $element.addClass('active');
        }
      };

    }
  };
}