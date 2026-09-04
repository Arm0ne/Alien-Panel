import{Di as e,Gn as t,Ii as n,Jr as r,Kn as i,On as a,Pi as o,Zn as s,jn as c,qr as l}from"./router-CdtSZDv0.js";import{n as u}from"./index-CD3NytCV.js";var d=l(`statistic`,[r(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),l(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[r(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[l(`icon`,{verticalAlign:`-0.125em`})]),r(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),r(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[l(`icon`,{verticalAlign:`-0.125em`})])])]),f=o({name:`Statistic`,props:Object.assign(Object.assign({},a.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(n){let{mergedClsPrefixRef:r,inlineThemeDisabled:o,mergedRtlRef:s}=i(n),l=a(`Statistic`,`-statistic`,d,u,n,r),f=c(`Statistic`,s,r),p=e(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:s},common:{cubicBezierEaseInOut:c}}=l.value;return{"--n-bezier":c,"--n-label-font-size":s,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=o?t(`statistic`,void 0,p,n):void 0;return{rtlEnabled:f,mergedClsPrefix:r,cssVars:o?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var e;let{mergedClsPrefix:t,$slots:{default:r,label:i,prefix:a,suffix:o}}=this;return(e=this.onRender)==null||e.call(this),n(`div`,{class:[`${t}-statistic`,this.themeClass,this.rtlEnabled&&`${t}-statistic--rtl`],style:this.cssVars},s(i,e=>n(`div`,{class:`${t}-statistic__label`},this.label||e)),n(`div`,{class:`${t}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},s(a,e=>e&&n(`span`,{class:`${t}-statistic-value__prefix`},e)),this.value===void 0?s(r,e=>e&&n(`span`,{class:`${t}-statistic-value__content`},e)):n(`span`,{class:`${t}-statistic-value__content`},this.value),s(o,e=>e&&n(`span`,{class:`${t}-statistic-value__suffix`},e))))}});export{f as t};