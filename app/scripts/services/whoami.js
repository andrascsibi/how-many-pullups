angular.module('pullApp')

.factory('WhoamiService', ['$q', '$http', '$route', function($q, $http, $route) {
  return function() {
    var delay = $q.defer();
    var promise = $http.get('/whoami');
    promise.success(function(data, status, headers, config) {
      delay.resolve(data);
    }).
    error(function(data, status, headers, config) {
      delay.reject('Could not fetch whoami');
    });
    return delay.promise;
  };
}]);
