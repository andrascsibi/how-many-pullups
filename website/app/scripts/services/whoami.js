angular.module('pullApp')

.factory('WhoamiService', ['$q', '$http', '$route', function($q, $http, $route) {
  return function() {
    var delay = $q.defer();
    var promise = $http.get('/whoami');
    promise.success(function(whoami, status, headers, config) {
      whoami.owner = $route.current.params.id === whoami.Account.ID;
      delay.resolve(whoami);
    }).
    error(function(data, status, headers, config) {
      delay.reject('Could not fetch whoami');
    });
    return delay.promise;
  };
}]);
