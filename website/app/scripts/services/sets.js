angular.module('pullApp')

.factory("Sets", ['$resource', function ($resource) {
  return $resource("/accounts/:id/challenges/:c_id/sets",
    {}, {});
}]);
