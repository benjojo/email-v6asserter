Email Delivery is stuck on IPv4
===

Generally speaking there is nothing that people want to talk about less than email delivery and for good reason, Email is continuously seen as one of those archaic protocols that everyone wants to improve but unfortunately has reached the point of critical mass where it is "unchangeable"

It was only recently when Google [started putting nasty icons on emails](https://www.blog.google/products/gmail/making-email-safer-for-you-posted-by/) when they were delivered to Gmail accounts unencrypted do people actually adopt TLS in a serious fashion.

One of those things that I've always wondered was that now that we have IPv6 in the world would it be possible to run an email server only on IPv6. Given that most service providers have IPv6 now surely they would be able to deliver email on v6?

To test this theory I figured it would be easy to write a script that would go through my inbox to see  what percentage of the emails inside were delivered over IPv6. I originally set to writing a [google appscript](https://developers.google.com/apps-script/) to do this function, unfortunately I realised that the limits on that system are pretty tight:

![google app script limits](/images/image2.png)

Fortunately the incredibly named "Google data Liberation Front" provides a tool called Google takeout that is capable of exporting your inbox into a single mbox file, With that I exported my traditional Gmail account and then my gsuite account emails for analysis.

![waiting for google takeout](/images/image4.png)

For the others wanted to do this it's worth knowing that the Google takeout tool does in fact take some serious time. My experience with the tool showed that it's exporting emails at less than one email per second. Something that is fine for small inboxes but if you have a large inbox then you may be spending a long time waiting for your zip file.

Then using the fantastic standard library of Google Go (mainly speaking, [net/mail](https://golang.org/pkg/net/mail/)), I made a small program that read the mbox and put out some basic CSV files with some basic statistics:

![first graph](/images/image5.gif)

As you can see IPv6 has always been very low and depending on your inbox diversity it has been getting lower, I Included the extra non google series in the graph above because on another Gmail account there was lots of emails from Google itself and it was inflating the IPv6 numbers quite significantly.

But there is hope! Surely if we managed to force email providers/senders to adopt TLS then we can get email providers to support IPv6 too?

Unfortunately not quite. TLS has been in use for much longer than IPv6 has, and in the case of at least my inbox emails that I am getting have been delivered by infrastructure as a service providers (IaaS) are going up, and these providers do not support IPv6 and don't seem like they are making any efforts to support IPv6 delivery.

![2nd graph](/images/image3.gif)

If just 4 IaaS email delivery providers were to support IPv6 then in the case of my inbox the IPv6 delivery count would go up by 28%

But this isn't the whole story. This is only looking at the inbox emails that are coming in, what if I was to run a email server that was only IPv6 and try to use it in the real world?

As I went down all the list of services it became immediately clear that almost none of them apart from Google services are actually able to send a confirmation email to a domain that only has an IPv6 mail server:

![google sheet screenshot](/images/image1.png)

I don't hold high hopes at any of this will be fixed in the near future, given the most of these sites themselves do not have IPv6 on their HTTP/HTTPS traffic then they have higher priorities than email delivery.

This is why I think that the responsibility shifts to these IaaS providers to support IPv6 Mail delivery. If 3 of these providers (Sendgrid, Mandrill, Amazon SES) were to support v6 they would lead the way for more systems to accept v6 mail delivery, giving that it would actually be in use.

Until then, I guess it’s just another waiting game.

As always, if you want to see the raw numbers, you can find them here: [Link](https://blog.benjojo.co.uk/asset/mJcMRsSlKF)

Or if you want to run the code yourself, you can find the program used here: [github](https://github.com/benjojo/email-v6asserter)