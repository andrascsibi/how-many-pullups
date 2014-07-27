angular.module('pullApp')

.factory("Stats", ['$resource', function ($resource) {
  return $resource("/accounts/:id/challenges/:c_id/stats",
    {}, {});
}]);
