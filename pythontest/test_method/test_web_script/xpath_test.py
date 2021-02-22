from lxml import etree

text = """
<div>
<ul>
<li class="item-o"><a href="link1.html">first item</a></li>
<li class="item-1"><a href="link2.html">second item</a></1i>
<li class="item-inactive"><a href="link3.html">third itemk</a></li>
<li class="item-1"><a href="link4.html">fourth item</a></li>
<li class="item-o"><a href="link5.html">fifth itemk</a>
</ul>
</div>
"""

# html = etree.HTML(text)
# result = etree.tostring(html)
# print(result.decode('utf-8'))

html = etree.parse('./index.html', etree.HTMLParser())
result = html.xpath('//*')
print(result)
