angular.module('pullApp')

.factory("Challenge", ['$resource', function ($resource) {
  return $resource("/accounts/:id/challenges/:c_id",
    {id: '@id', c_id: '@c_id'}, {});
}]);
