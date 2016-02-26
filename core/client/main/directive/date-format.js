angular
  .module('app')
  .directive('dateFormat', directiveDateFormat);

function directiveDateFormat() {
  return {
    restrict: 'A',
    require: "?ngModel",
    link: function(scope, element, attrs, ngModel) {
      var dateFormat;

      if (!ngModel) return;

      attrs.$observe('dateFormat', function(value) {
        dateFormat = value;
        ngModel.$render();
      });

      ngModel.$render = function() {
        element.val(ngModel.$modelValue ? moment(ngModel.$modelValue).format(dateFormat) : undefined);
        scope.ngModel = ngModel.$modelValue;
      };

      ngModel.$parsers.unshift(function(viewValue) {
        var date = moment(viewValue);
        ngModel.$setValidity('date', !viewValue || date.isValid());
        return date.isValid()? date.format(dateFormat): moment().format(dateFormat);
      });

      // format the new model value before it is displayed
      ngModel.$formatters.push(function(modelValue) {
        var date = moment(modelValue);
        ngModel.$setValidity('date', !modelValue || date.isValid());

        return moment(modelValue).format(dateFormat);
      });
    }
  };
}
