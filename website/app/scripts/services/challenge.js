angular.module('pullApp')

.factory("Challenge", ['$resource', function ($resource) {
  return $resource("/accounts/:id/challenges/:c_id",
    {id: '@AccountID', c_id: '@ID'}, {});
}]);
