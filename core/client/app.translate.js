angular
  .module('app')
  .config(appTranslate);

appTranslate.$inject = ['$translateProvider'];
function appTranslate ($translateProvider) {
  $translateProvider.useStaticFilesLoader({
    prefix: './assets/i18n/',
    suffix: '.json'
  }).fallbackLanguage('en');
}