window.Todos = Ember.Application.create();

Todos.ApplicationAdapter = DS.RESTAdapter.extend({
  namespace: 'api',
  host: 'http://localhost:8080',
  buildURL: function (type,id,record) {
    url = this._super(type,"",record)
    // This will replace abc.com/todos/{id} with abc.com/todos?id={id}
    if(id !== undefined)
      return url + "?id=" + id
    return url
  },
});

Todos.DateTransform = DS.DateTransform.extend({
  deserialize: function(serialized) {
    return serialized ? null : new Date(deserialized);
  },
  serialize: function(deserialized) {
    return Ember.isEmpty(deserialized) ? null : (new Date(deserialized)).toISOString();
  }
})