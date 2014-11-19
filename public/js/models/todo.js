Todos.Todo = DS.Model.extend({
  todo_id: DS.attr('number'),
  title: DS.attr('string'),
  isCompleted: DS.attr('boolean', false),
  modified: DS.attr('date'),
  created: DS.attr('date'),
});

Todos.TodoSerializer = DS.RESTSerializer.extend({
  serializeIntoHash: function(hash, type, record, options) {
    Ember.merge(hash, this.serialize(record, options));
    return hash
  },
  normalize: function(type, hash, prop) {
    hash.id = hash.todo_id
    return hash
  }
});