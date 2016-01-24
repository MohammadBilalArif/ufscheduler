# UFScheduler - GatorPik

The GatorPik web application is a tool for University of Florida students to determine the order of classes to take in future semesters.  By copy+pasting their Degree Audit
into the web app, we're able to scrape their completed classes and class requirements.  We then build a dependency graph using real-time data scraped from the University of
Florida Course Listing to determine which classes are prerequisites for others and the students class requirements to build an ordered list of classes to take first at the
front and classes to take last at the end.  We display this information to the student using a web interface that includes the feature of highlighting prerequisites for
classes on mouse over.

# Instructions for Use

Visit www.gatorpik.com and follow the instructions on screen, summarized here for reference:

  1) Visit ISIS - the University of Florida Student Information  System

  2) Select Degree Audit on the left.

  3) Log in with your student information.

  4) Click on View Critical Tracking Audit / View Unmet Requirements as required.

  5) Select all (ctrl+a) then copy (ctrl+c).

  6) Paste (ctrl+v) the information into the appropriate box.

By allowing the student to copy and paste their entire Degree Audit into the webpage, we attempted to simplify the user's interaction with the app.

# How to Read the Results

When your class list has been calculated, a list of classes will appear, with classes you should take first at the top and those to take later at the bottom.  *Be warned!* The
system currently does not distinguish between optional prerequisites and required ones, as well as technical electives, so you should use the list as a **reference** in taking
the indicated classes, not as the only possible route!

# License

This code is released under the MIT license.

Copyright (c) 2016 Chris Pergrossi and Eric McGhee.


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
