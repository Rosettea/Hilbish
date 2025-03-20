---
title: Module hilbish.messages
description: simplistic message passing
layout: doc
menu:
  docs:
    parent: "API"
---


## Introduction
The messages interface defines a way for Hilbish-integrated commands,
user config and other tasks to send notifications to alert the user.z
The `hilbish.message` type is a table with the following keys:
`title` (string): A title for the message notification.
`text` (string): The contents of the message.
`channel` (string): States the origin of the message, `hilbish.*` is reserved for Hilbish tasks.
`summary` (string): A short summary of the `text`.
`icon` (string): Unicode (preferably standard emoji) icon for the message notification
`read` (boolean): Whether the full message has been read or not.

## Functions
|||
|----|----|
|<a href="#unreadCount">unreadCount()</a>|Returns the amount of unread messages.|
|<a href="#readAll">readAll()</a>|Marks all messages as read.|
|<a href="#send">send(message)</a>|Sends a message.|
|<a href="#read">read(idx)</a>|Marks a message at `idx` as read.|
|<a href="#delete">delete(idx)</a>|Deletes the message at `idx`.|
|<a href="#clear">clear()</a>|Deletes all messages.|
|<a href="#all">all()</a>|Returns all messages.|
<hr>
<div id='all'>
<h4 class='heading'>
hilbish.messages.all()
<a href="#all" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns all messages.
#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='clear'>
<h4 class='heading'>
hilbish.messages.clear()
<a href="#clear" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Deletes all messages.
#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='delete'>
<h4 class='heading'>
hilbish.messages.delete(idx)
<a href="#delete" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Deletes the message at `idx`.
#### Parameters
`idx` **`number`**  


</div>

<hr>
<div id='read'>
<h4 class='heading'>
hilbish.messages.read(idx)
<a href="#read" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Marks a message at `idx` as read.
#### Parameters
`idx` **`number`**  


</div>

<hr>
<div id='send'>
<h4 class='heading'>
hilbish.messages.send(message)
<a href="#send" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Sends a message.
#### Parameters
`message` **`hilbish.message`**  


</div>

<hr>
<div id='readAll'>
<h4 class='heading'>
hilbish.messages.readAll()
<a href="#readAll" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Marks all messages as read.
#### Parameters
This function has no parameters.  
</div>

<hr>
<div id='unreadCount'>
<h4 class='heading'>
hilbish.messages.unreadCount()
<a href="#unreadCount" class='heading-link'>
	<i class="fas fa-paperclip"></i>
</a>
</h4>

Returns the amount of unread messages.
#### Parameters
This function has no parameters.  
</div>

