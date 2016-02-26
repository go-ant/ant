angular
  .module('app')
  .config(appConfig);

appConfig.$inject = ['flowFactoryProvider', 'localStorageServiceProvider'];
function appConfig (flowFactoryProvider, localStorageServiceProvider) {
  flowFactoryProvider.defaults = {
    target: '../api/backend/upload',
    fileParameterName: 'file',
    chunkSize: 5<<20,
    singleFile: true,
    testMethod: 'Post',
    testChunks: false
  };

  localStorageServiceProvider.setPrefix('goant');
}