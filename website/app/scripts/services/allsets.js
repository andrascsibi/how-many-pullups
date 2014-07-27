angular.module('pullApp')

.factory("AllSets", ['$resource', function ($resource) {
  return $resource("/accounts/:id/challenges/:c_id/export",
    {}, {});
}]);
